package remove

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
)

var iterator_uri string
var exporter_uri string
var writer_uri string
var verbose bool

var properties multi.MultiString

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("remove")

	fs.StringVar(&iterator_uri, "iterator-uri", "repo://", "A valid whosonfirst/go-whosonfirst-iterate/v3.Iterator URI.")
	fs.StringVar(&exporter_uri, "exporter-uri", "whosonfirst://", "A valid whosonfirst/go-whosonfirst-export/v3.Exporter URI.")
	fs.StringVar(&writer_uri, "writer-uri", "null://", "A valid whosonfirst/go-writer/v3.Writer URI.")
	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging")

	fs.Var(&properties, "property", "One or more (fully-qualified) paths for a property to remove.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Remove properties from one or more Who's On First records.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t%s [options] path(N) path(N)\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n")
		fs.PrintDefaults()
	}

	return fs
}
