package srv

import (
	"context"
	"lru-cache/internal/cache"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

type Server struct {
	storage cache.ILRUCache
	router  chi.Router
	cfg     Config
}

type Config struct {
	HostPort   string
	Capacity   int
	DefaultTTL time.Duration
}

func New(cfg Config) Server {
	storage := cache.New(cfg.Capacity)
	router := chi.NewRouter()

	return Server{storage: storage, router: router, cfg: cfg}
}

func (s *Server) Run(ctx context.Context) error {

	s.router.Route("/api/lru", func(r chi.Router) {
		r.Post("/", s.putKey)
		r.Get("/{key}", s.getKey)
		r.Get("/", s.getAllKeys)
		r.Delete("/{key}", s.evictKey)
		r.Delete("/", s.evictAllKeys)
	})

	err := http.ListenAndServe(s.cfg.HostPort, s.router)
	if err != nil {
		return err
	}

	return nil
}
