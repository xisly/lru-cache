package srv

import (
	"log/slog"
	"net/http"
	"time"

	"lru-cache/internal/models"
	"lru-cache/pkg/errs"

	"github.com/go-chi/chi/v5"
)

func (s *Server) postKey(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data := &models.PostRequest{}
	err := data.FromJSON(r.Body)
	if err != nil {
		s.logger.Debug("Unable to unmarshall JSON in post key", slog.Any("error", err))
		http.Error(rw, "Invalid request body", http.StatusBadRequest)
		return
	}

	ttl := s.cfg.DefaultTTL
	if data.TTLSeconds > 0 {
		ttl = time.Duration(data.TTLSeconds) * time.Second
	}

	err = s.storage.Put(ctx, data.Key, data.Value, ttl)
	if err != nil {
		s.logger.Warn("Something went wrong in post a key", slog.Any("error", err))
		http.Error(rw, "Something went wrong", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)

	s.logger.Debug("Created a key", slog.String("key", data.Key), slog.Any("value", data.Value), slog.Duration("ttl", ttl))
}

func (s *Server) getKey(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	key := chi.URLParam(r, "key")
	data := &models.GetResponse{Key: key}
	var err error
	data.Value, data.TimeExpiresAt, err = s.storage.Get(ctx, key)
	if err != nil {
		if err == errs.ErrNotFound {
			s.logger.Debug("Key not found in get by key", slog.String("key", key))
			http.Error(rw, "Not found", http.StatusNotFound)
		} else {
			s.logger.Warn("Something went wrong in get by key", slog.String("key", key))
			http.Error(rw, "Something went wrong", http.StatusInternalServerError)
		}
		return
	}

	rw.WriteHeader(http.StatusOK)
	if err := data.ToJSON(rw); err != nil {
		s.logger.Error("Failed to marshall struct into JSON", slog.String("key", key))
		http.Error(rw, "Something went horribly wrong", http.StatusInternalServerError)
	}

	s.logger.Debug("Got data by key", slog.String("key",data.Key), slog.Any("value", data.Value), slog.Time("expires at", data.TimeExpiresAt))
}

func (s *Server) getAllKeys(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data := &models.GetAllResponse{}

	var err error
	data.Keys, data.Values, err = s.storage.GetAll(ctx)
	if err != nil {
		if err == errs.ErrCacheIsEmpty {
			s.logger.Debug("Cache is empty, unable to get all keys")
			rw.WriteHeader(http.StatusNoContent)
		} else {
			s.logger.Warn("Something went wrong in get all keys", slog.Any("Error:", err))
			http.Error(rw, "Something went wrong", http.StatusInternalServerError)
		}
		return
	}

	rw.WriteHeader(http.StatusOK)
	if err := data.ToJSON(rw); err != nil {
		s.logger.Warn("Unable to marshall data into JSON in get all keys")
		http.Error(rw, "Something went wrong", http.StatusInternalServerError)
	}

	s.logger.Debug("Got all keys")
}

func (s *Server) evictKey(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	key := chi.URLParamFromCtx(ctx, "key")

	value, err := s.storage.Evict(ctx, key)
	if err != nil {
		if err == errs.ErrNotFound {
			s.logger.Debug("Key not found in delete by key", slog.String("key", key))
			http.Error(rw, "Not found", http.StatusNotFound)
		} else {
			s.logger.Warn("Something went wrong in delete by key", slog.String("key", key))
			http.Error(rw, "Something went wrong", http.StatusInternalServerError)
		}
		return
	}
	rw.WriteHeader(http.StatusNoContent)

	s.logger.Debug("Deleted a key", slog.String("key", key), slog.Any("value", value))
}

func (s *Server) evictAllKeys(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err := s.storage.EvictAll(ctx)
	if err != nil {
		if err == errs.ErrCacheIsEmpty {
			s.logger.Debug("No keys deleted due to cache emptyness")
			http.Error(rw, "Cache is empty", http.StatusNoContent)
		} else {
			s.logger.Warn("Something went horribly wrong in delete all keys", slog.Any("error", err))
			http.Error(rw, "Something went wrong", http.StatusInternalServerError)
		}
		return
	}
	rw.WriteHeader(http.StatusNoContent)

	s.logger.Debug("Deleted all keys")
}
