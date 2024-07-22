package srv

import (
	"context"
	"errors"
	"log/slog"
	"lru-cache/internal/cache"
	"lru-cache/pkg/loghandler"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	storage cache.ILRUCache
	router  chi.Router
	cfg     Config
	logger  *slog.Logger
}

type Config struct {
	HostPort   string        `env:"SERVER_HOST_PORT" envDefault:"localhost:8080"`
	CacheSize  int           `env:"CACHE_SIZE" envDefault:"10"`
	DefaultTTL time.Duration `env:"DEFAULT_CACHE_TTL" envDefault:"1m"`
	LogLevel   string        `env:"LOG_LEVEL" envDefault:"DEBUG"`
}

func New(cfg Config) (*Server, error) {
	levelhandler, err := loghandler.New(cfg.LogLevel, slog.NewTextHandler(os.Stdout, nil))
	if err != nil {
		return nil, err
	}

	logger := slog.New(levelhandler)

	storage := cache.New(cfg.CacheSize)

	router := chi.NewRouter()

	return &Server{storage: storage, router: router, cfg: cfg, logger: logger}, nil
}

func (s *Server) Run(ctx context.Context) error {

	s.router.Use(s.loggingMiddleware, middleware.Recoverer)

	s.router.Route("/api/lru", func(r chi.Router) {
		r.Post("/", s.putKey)
		r.Get("/{key}", s.getKey)
		r.Get("/", s.getAllKeys)
		r.Delete("/{key}", s.evictKey)
		r.Delete("/", s.evictAllKeys)
	})

	server := http.Server{
		Addr: s.cfg.HostPort,
		Handler: s.router,
	}

	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("HTTP server error", slog.Any("error", err))
			os.Exit(1)
		}
		s.logger.Info("Stopped serving new connections.")
	}()
	
	s.logger.Info("Running server")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	
	shutdownCtx, shutdownRelease := context.WithTimeout(ctx, 10*time.Second)
	defer shutdownRelease()
	
	if err := server.Shutdown(shutdownCtx); err != nil {
		s.logger.Error("HTTP shutdown error", slog.Any("error",err))
		os.Exit(2)
	}
	s.logger.Info("Graceful shutdown complete.")

	return nil
}