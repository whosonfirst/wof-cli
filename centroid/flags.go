package centroid

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
)

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("centroid")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Emit centroids data (source, latitude, longitude) for one or more Who's On First records as CSV-encoded properties.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s emit [options] path(N) path(N)\n", os.Args[0])
		fs.PrintDefaults()
	}

	return fs
}
