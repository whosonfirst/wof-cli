//go:build uri

package wof

// Experimental: Moving specific imports in to build tags so people to build trimmed-down
// binaries if they want to

import (
	_ "github.com/whosonfirst/wof/uri"
)
