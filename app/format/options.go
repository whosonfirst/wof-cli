package format

import (
	"context"
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	URIs []string
	Overwrite bool
}

func RunOptionsFromFlagSet(ctx context.Context, fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	uris := fs.Args()

	opts := &RunOptions{
		URIs: uris,
	}

	return opts, nil
}
