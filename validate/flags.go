package validate

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("validate")

	return fs
}
