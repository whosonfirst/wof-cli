package export

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
)

var stdout bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("export")
	fs.BoolVar(&stdout, "stdout", true, "Boolean flag signaling that updated records should be written to STDOUT. If false input files will be overwritten.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "\"Export-ify\" one or more Who's On First records.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s path(N) path(N)\n", os.Args[0])
		fs.PrintDefaults()
	}

	return fs
}
