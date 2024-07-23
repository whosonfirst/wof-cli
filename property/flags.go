package property

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
)

var paths multi.MultiString
var format string
var prefix string
var action string

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("properties")
	fs.Var(&paths, "path", "One or more valid tidwall/gjson paths to extract from each document")
	fs.StringVar(&format, "format", "", "Valid options are: csv. If empty then properties will printed as a new-line separated list.")
	fs.StringVar(&prefix, "prefix", "", "If not empty this prefix will be appended (and separated by a \".\") to each -path argument")
	fs.StringVar(&action, "action", "list", "Valid options are: list, remove, replace.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Print one or more properties for one or more Who's On First IDs.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s path(N) path(N)\n", os.Args[0])
		fs.PrintDefaults()
	}

	return fs
}
