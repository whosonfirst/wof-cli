//go:build emit_git_iterator

package emit

// Experimental: Moving specific imports in to build tags so people to build trimmed-down
// binaries if they want to

import (
	_ "github.com/whosonfirst/go-whosonfirst-iterate-git/v2"
)
