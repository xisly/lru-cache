package loglevel

import (
	"log/slog"
	"lru-cache/pkg/errs"
)

func Set(lv string) (slog.Level, error) {
	switch lv {
	case "WARN":
		return slog.LevelWarn, nil
	case "DEBUG":
		return slog.LevelDebug, nil
	case "INFO":
		return slog.LevelInfo, nil
	case "ERROR":
		return slog.LevelError, nil
	default:
		return slog.LevelDebug, errs.ErrIncorrectLogLevel
	}

}
