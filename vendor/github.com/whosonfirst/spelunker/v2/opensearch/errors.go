package opensearch

import (
	"errors"
)

// ErrCursorIsExpired returns an error signaling that an OpenSearch cursor has expired.
var ErrCursorIsExpired = errors.New("Query cursor has expired")
