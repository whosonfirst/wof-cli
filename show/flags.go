package show

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
)

var port int

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("validate")
	fs.IntVar(&port, "port", 8080, "...")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Display one or more Who's On First documents on a map.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s path(N) path(N)\n", os.Args[0])
		fs.PrintDefaults()
	}

	return fs
}