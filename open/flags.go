package open

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
)

var editor string

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("validate")
	fs.StringVar(&editor, "editor", "", "The editor to use for opening Who's On First records. If empty the value of the EDITOR environment variable will be used.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Open one or more Who's On First documents in a custom editor.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s path(N) path(N)\n", os.Args[0])
		fs.PrintDefaults()
	}

	return fs
}
