//go:build !no_spatial_sqlite

package pip

// Experimental: Moving specific imports in to build tags so people to build trimmed-down
// binaries if they want to

import (
	_ "github.com/whosonfirst/go-whosonfirst-spatial-sqlite"
)
