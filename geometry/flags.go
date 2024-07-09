package geometry

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
)

var action string
var source string
var format string
var stdout bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("emit")

	fs.StringVar(&action, "action", "show", "Valid options are: show, update.")
	fs.StringVar(&format, "format", "geojson", "The format of the source geometry (if -action update). Valid options are: geojson.")
	fs.StringVar(&source, "source", "", "The data source from which to derive a geometry (if -action update). Valid options are: the path to a local file.")
	fs.BoolVar(&stdout, "stdout", false, "Boolean flag signaling that updated records should be written to STDOUT. If false input files will be overwritten.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Display or update the geometry for one or more Who's On First records.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s geometry [options] path(N) path(N)\n", os.Args[0])
		fs.PrintDefaults()
	}

	return fs
}
