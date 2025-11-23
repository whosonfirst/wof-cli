package spelunker

import (
	"errors"
)

// ErrNotImplemented returns an error signaling a method or feature is not implemented.
var ErrNotImplemented = errors.New("Not implemented")

// ErrNotFound returns an error signaling a record has not been indexed or is not present.
var ErrNotFound = errors.New("Not found")
