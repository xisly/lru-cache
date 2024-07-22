//Package levelhandler provides a slog handler that manages filtering of logs based on
// a preconfigured level.
package levelhandler

import (
	"context"
	"log/slog"
	"lru-cache/pkg/errs"
)

// LevelHandler wraps another handler and filters logs based on the level.
type LevelHandler struct {
	level   slog.Level
	handler slog.Handler
}

// New returns a LevelHandler with the given level.
// All methods except Enabled delegate to h.
func New(level string, handler slog.Handler) (slog.Handler, error) {
	lvl, err := parseLevel(level)
	if err != nil {
		return nil, err
	}
	return &LevelHandler{
		level:   lvl,
		handler: handler,
	}, nil
}

// Enabled checks if the log level is enabled.
func (h *LevelHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level
}

// Handle handles the log record if the level is enabled.
func (h *LevelHandler) Handle(ctx context.Context, r slog.Record) error {
	if h.Enabled(ctx, r.Level) {
		return h.handler.Handle(ctx, r)
	}
	return nil
}

// WithAttrs returns a new LevelHandler with the given attributes.
func (h *LevelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &LevelHandler{level: h.level, handler: h.handler.WithAttrs(attrs)}
}

// WithGroup returns a new LevelHandler with the given group.
func (h *LevelHandler) WithGroup(name string) slog.Handler {
	return &LevelHandler{level: h.level, handler: h.handler.WithGroup(name)}
}

func parseLevel(level string) (slog.Level, error) {
	switch level {
	case "WARN":
		return slog.LevelWarn, nil
	case "DEBUG":
		return slog.LevelDebug, nil
	case "INFO":
		return slog.LevelInfo, nil
	case "ERROR":
		return slog.LevelError, nil
	}

	return 0, errs.ErrIncorrectLogLevel
}
