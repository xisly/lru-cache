package errs

import "errors"

var (
	ErrNotFound          = errors.New("value not found")
	ErrCacheIsEmpty      = errors.New("cache is empty")
	ErrIncorrectLogLevel = errors.New("Unable to parse log level")
)
