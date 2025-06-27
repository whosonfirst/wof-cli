//go:build !no_iterateor_org

package emit

// Experimental: Moving specific imports in to build tags so people to build trimmed-down
// binaries if they want to

import (
	_ "github.com/whosonfirst/go-whosonfirst-iterate-git/v3/github"
)
