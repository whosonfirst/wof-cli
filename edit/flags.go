package edit

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
)

var reader_uri string
var writer_uri string

var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("edit")

	fs.StringVar(&reader_uri, "reader-uri", "fs:///", "...")
	fs.StringVar(&writer_uri, "writer-uri", "fs:///", "...")
	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Launch a local web server running on a random port hosting a web application for editing one or more Who's On First records.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s path(N) path(N)\n", os.Args[0])
		fs.PrintDefaults()
	}

	return fs
}
