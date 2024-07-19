package srv

import (
	"github.com/go-chi/chi"
	"lru-cache/internal/cache"
	"time"
)

type Server struct {
	storage  cache.ILRUCache
	router chi.Router
	cfg    Config
}

type Config struct {
	HostPort   string
	Capacity   int
	DefaultTTL time.Duration
}

func New(cfg Config) Server {
	 storage := cache.New(cfg.Capacity)
   router := chi.NewRouter()

   return Server {storage: storage, router: router, cfg: cfg}
}

func (s *Server) Run(ctx context.Context) error {

  s.router.Route("/api/lru", func (r chi.Router) {
    r.Post("/", s.put)
    r.Get("/{key}", s.get)
    r.Get("/", s.getAll)
    r.Delete("/{key}",evict)
    r.Delete("/",evictAll)
  })

	http.ListenAndServe(s.cfg.HostPort, s.router)
}
