package ensure

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
)

var iterator_uri string
var exporter_uri string
var writer_uri string
var verbose bool

var str_properties multi.KeyValueString
var int_properties multi.KeyValueInt64
var float_properties multi.KeyValueFloat64
var bool_properties multi.KeyValueBool
var geom_property multi.KeyValueStringFlag

var if_missing bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("ensure")

	fs.StringVar(&iterator_uri, "iterator-uri", "repo://", "A valid whosonfirst/go-whosonfirst-iterate/v3.Iterator URI.")
	fs.StringVar(&exporter_uri, "exporter-uri", "whosonfirst://", "A valid whosonfirst/go-whosonfirst-export/v3.Exporter URI.")
	fs.StringVar(&writer_uri, "writer-uri", "null://", "A valid whosonfirst/go-writer/v3.Writer URI.")
	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging")

	fs.Var(&str_properties, "string-property", "One or more {KEY}={VALUE} flag where {KEY} is a valid tidwall/gjson path and {VALUE} is a string value.")
	fs.Var(&int_properties, "int-property", "One or more {KEY}={VALUE} flag where {KEY} is a valid tidwall/gjson path and {VALUE} is a int(64) value.")
	fs.Var(&float_properties, "float-property", "One or more {KEY}={VALUE} flags where {KEY} is a valid tidwall/gjson path and {VALUE} is a float(64) value.")
	fs.Var(&bool_properties, "boolean-property", "One or more {KEY}={VALUE} flags where {KEY} is a valid tidwall/gjson path and {VALUE} is a boolean value.")
	fs.Var(&geom_property, "geometry-property", "A {KEY}={VALUE} flag indicating the source of the geometry data to assign. Valid options are: wkt={VALID_WKT_GEOMETRY}, geojson={VALID_GEOJSON_GEOMETRY}, file={PATH_TO_GEOJSON_FEATURE}.")

	fs.BoolVar(&if_missing, "if-missing", false, "Only assign property value if the property key is not set.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Ensure typed values for one or more Who's On First records.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t%s [options] path(N) path(N)\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n")
		fs.PrintDefaults()
	}

	return fs
}
