package supersede

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
)

var superseding_id int64
var parent_id int64

var is_deprecated bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("properties")
	fs.Int64Var(&parent_id, "parent_id", -1, "The parent ID to assign to the new record.")
	fs.Int64Var(&superseding_id, "superseding_id", -1, "The ID to supersede each record with. If `-1` then each record will be cloned and the new ID of the clone will be used as the superseding_id.")

	fs.BoolVar(&is_deprecated, "deprecated", false, "Each record being superseded should also be marked as deprecated.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Print one or more properties for one or more Who's On First IDs.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s path(N) path(N)\n", os.Args[0])
		fs.PrintDefaults()
	}

	return fs
}
