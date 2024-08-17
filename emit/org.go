//go:build emit_org_iterator

package emit

// Experimental: Moving specific imports in to build tags so people to build trimmed-down
// binaries if they want to

import (
	_ "github.com/whosonfirst/go-whosonfirst-iterate-organization/v2"
)
