//go:build emit_geoparquet_writer

package emit

// Experimental: Moving specific imports in to build tags so people to build trimmed-down
// binaries if they want to

import (
	_ "github.com/whosonfirst/go-writer-geoparquet/v3"
)
