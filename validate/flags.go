package validate

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
)

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("validate")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Validate one or more Who's On First documents.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s path(N) path(N)\n", os.Args[0])
		fs.PrintDefaults()
	}

	return fs
}
