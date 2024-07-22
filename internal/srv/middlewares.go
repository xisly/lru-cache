package srv

import (
	"log/slog"
	"net/http"
	"time"
)

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r.WithContext(r.Context()))

		s.logger.Debug("handled request",slog.Time("time", time.Now()),slog.String("method",r.Method),slog.String("URI", r.RequestURI), slog.Duration("handling time", time.Since(start)))
	})
}
