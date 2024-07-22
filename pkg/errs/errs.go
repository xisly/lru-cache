//Package errs provides common errors handled in the project.
package errs

import "errors"

var (
	//ErrNotFound is used when a value's not found in cache.
	ErrNotFound          = errors.New("value not found")
	//ErrCacheIsEmpty is used when it's impossible to get all the cache 
	//or delete the entire cache due to it's emptyness.
	ErrCacheIsEmpty      = errors.New("cache is empty")
	//ErrIncorrectLogLevel is used when it's impossible to parse log level
	//from a flag or env.
	ErrIncorrectLogLevel = errors.New("unable to parse log level")
)
