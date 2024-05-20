package pip

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
)

var mapshaper_client_uri string
var spatial_database_uri string
var placetype string

var stdout bool
var exportify bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("pip")
	fs.StringVar(&mapshaper_client_uri, "mapshaper-client-uri", "", "Optional URI to a sfomuseum/go-sfomuseum-mapshaper server instance used to derive point-in-polygon centroids. If absent then the centroid used to perform point-in-polygon operations will be determined using internal heuristics.")
	fs.StringVar(&spatial_database_uri, "spatial-database-uri", "", "A valid whosonfirst/go-whosonfirst-spatial/database URI. By default 'pmtiles://' and 'sqlite://' spatial database URIs are supported.")
	fs.StringVar(&placetype, "placetype", "", "Assign this value as the \"wof:placetype\" property before performing point-in-polygon operations.")
	fs.BoolVar(&stdout, "stdout", true, "Boolean flag signaling that updated records should be written to STDOUT.")
	fs.BoolVar(&exportify, "export", false, "\"Export-ify\" each record after point-in-polygon operations are complete.")

	// To do: Add input filter flags

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Perform point-in-polygon and wof:hierarchy update operations on one or more Who's On First records.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s path(N) path(N)\n", os.Args[0])
		fs.PrintDefaults()
	}

	return fs
}
