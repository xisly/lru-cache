// Package srv provides a configured cache server with graceful shudown logic a function to run it,
// along with all the http handlers, middlewares and test used in a project.
package srv

import (
	"context"
	"errors"
	"log/slog"
	"lru-cache/internal/cache"
	"lru-cache/pkg/levelhandler"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Server defines a configured server with storage, router, configuration, and logger.
type Server struct {
	storage cache.ILRUCache
	router  chi.Router
	cfg     Config
	logger  *slog.Logger
}

// Config holds the configuration parameters for the server.
type Config struct {
	HostPort   string        `env:"SERVER_HOST_PORT" envDefault:"localhost:8080"`
	CacheSize  int           `env:"CACHE_SIZE" envDefault:"10"`
	DefaultTTL time.Duration `env:"DEFAULT_CACHE_TTL" envDefault:"1m"`
	LogLevel   string        `env:"LOG_LEVEL" envDefault:"WARN"`
}

// New creates a new Server with the provided configuration.
// It initializes the logger, storage, and router.
// Returns a configured Server instance or an error if initialization fails.
func New(cfg Config) (*Server, error) {
	levelhandler, err := levelhandler.New(cfg.LogLevel, slog.NewTextHandler(os.Stdout, nil))
	if err != nil {
		return nil, err
	}

	logger := slog.New(levelhandler)

	storage := cache.New(cfg.CacheSize)
	logger.Info("Created LRU cache", slog.Int("size", cfg.CacheSize))

	router := chi.NewRouter()

	logger.Debug("Configured", slog.Any("config", cfg))

	return &Server{storage: storage, router: router, cfg: cfg, logger: logger}, nil
}

// Run starts the server and listens for incoming HTTP requests.
// It sets up the routes and middleware, and handles graceful shutdown on receiving a termination signal.
// Returns an error if the server encounters issues during operation.
func (s *Server) Run(ctx context.Context) error {

	s.router.Use(s.loggingMiddleware, middleware.Recoverer)

	s.router.Route("/api/lru", func(r chi.Router) {
		r.Post("/", s.postKey)
		r.Get("/{key}", s.getKey)
		r.Get("/", s.getAllKeys)
		r.Delete("/{key}", s.evictKey)
		r.Delete("/", s.evictAllKeys)
	})

	server := http.Server{
		Addr:    s.cfg.HostPort,
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
		s.logger.Error("HTTP shutdown error", slog.Any("error", err))
		os.Exit(2)
	}
	s.logger.Info("Graceful shutdown complete.")

	return nil
}