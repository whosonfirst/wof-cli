package uri

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
)

var prefix string

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("validate")
	fs.StringVar(&prefix, "prefix", "", "An optional prefix to append to the final URI.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Print the nested URI for one or more Who's On First IDs.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s path(N) path(N)\n", os.Args[0])
		fs.PrintDefaults()
	}

	return fs
}
