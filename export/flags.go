package export

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var overwrite bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("export")
	fs.BoolVar(&overwrite, "overwrite", false, "Boolean flag signaling that the source file should be overwritten.")

	return fs
}
