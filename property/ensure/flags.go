package ensure

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
)

var iterator_uri string
var exporter_uri string
var writer_uri string

var str_properties multi.KeyValueString
var int_properties multi.KeyValueInt64
var float_properties multi.KeyValueFloat64

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("ensure")

	fs.StringVar(&iterator_uri, "iterator-uri", "repo://", "A valid whosonfirst/go-whosonfirst-iterate/v3.Iterator URI.")
	fs.StringVar(&exporter_uri, "exporter-uri", "whosonfirst://", "A valid whosonfirst/go-whosonfirst-export/v3.Exporter URI.")
	fs.StringVar(&writer_uri, "writer-uri", "null://", "A valid whosonfirst/go-writer/v3.Writer URI.")

	fs.Var(&str_properties, "string-property", "One or more {KEY}={VALUE} flag where {KEY} is a valid tidwall/gjson path and {VALUE} is a string value.")

	fs.Var(&int_properties, "int-property", "One or more {KEY}={VALUE} flag where {KEY} is a valid tidwall/gjson path and {VALUE} is a int(64) value.")

	fs.Var(&float_properties, "float-property", "One or more {KEY}={VALUE} flags where {KEY} is a valid tidwall/gjson path and {VALUE} is a float(64) value.")

	return fs
}
