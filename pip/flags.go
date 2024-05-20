package pip

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var mapshaper_client_uri string
var spatial_database_uri string

var overwrite bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("pip")
	fs.StringVar(&mapshaper_client_uri, "mapshaper-client-uri", "", "Optional URI to a sfomuseum/go-sfomuseum-mapshaper server instance used to derive point-in-polygon centroids. If absent then the centroid used to perform point-in-polygon operations will be determined using internal heuristics.")
	fs.StringVar(&spatial_database_uri, "spatial-database-uri", "", "A valid whosonfirst/go-whosonfirst-spatial/database URI. By default 'pmtiles://' and 'sqlite://' spatial database URIs are supported.")
	fs.BoolVar(&overwrite, "overwrite", false, "Boolean flag signaling that the source file should be overwritten.")

	return fs
}
