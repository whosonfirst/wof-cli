package edit

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
)

var writer_uri string

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("edit")
	fs.StringVar(&writer_uri, "writer-uri", "stdout://", "...")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "\"Export-ify\" one or more Who's On First records.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s path(N) path(N)\n", os.Args[0])
		fs.PrintDefaults()
	}

	return fs
}
