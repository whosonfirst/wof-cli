package rebuild

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
var parent_reader_uri string
var verbose bool

var parent_ids multi.MultiInt64

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("remove")

	fs.StringVar(&iterator_uri, "iterator-uri", "repo://", "A valid whosonfirst/go-whosonfirst-iterate/v3.Iterator URI.")
	fs.StringVar(&exporter_uri, "exporter-uri", "whosonfirst://", "A valid whosonfirst/go-whosonfirst-export/v3.Exporter URI.")
	fs.StringVar(&writer_uri, "writer-uri", "null://", "A valid whosonfirst/go-writer/v3.Writer URI.")
	fs.StringVar(&parent_reader_uri, "parent-reader-uri", "https://data.whosonfirst.org", "A valid whosonfirst/go-reader/v2.Reader URI.")
	fs.Var(&parent_ids, "parent-id", "One or more explicit parent IDs to use for deriving hierarchies. The default is to use each record's wof:parent_id value.")
	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Rebuild the wof:hierarchy property one or more Who's On First records.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t%s [options] path(N) path(N)\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n")
		fs.PrintDefaults()
	}

	return fs
}
