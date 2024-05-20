package format

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
)

var stdout bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("format")
	fs.BoolVar(&stdout, "stdout", true, "Boolean flag signaling that updated records should be written to STDOUT. If false input files will be overwritten.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Format one or more GeoJSON files according to Who's On First conventions.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s path(N) path(N)\n", os.Args[0])
		fs.PrintDefaults()
	}

	return fs
}
