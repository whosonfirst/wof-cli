package emit

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
)

var writer_uri string
var iterator_uri string

var forgiving bool

var as_spr bool
var as_spr_geojson bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("emit")

	fs.StringVar(&writer_uri, "writer-uri", "jsonl://?writer=stdout://", "A valid whosonfirst/go-writer.Writer URI.")
	fs.StringVar(&iterator_uri, "iterator-uri", "repo://", "A valid whosonfirst/go-whosonfirst-iterate/v2/emitter URI.")

	fs.BoolVar(&as_spr, "as-spr", false, "Emit Who's On First records formatted as Standard Place Response (SPR) records.")
	fs.BoolVar(&as_spr_geojson, "as-spr-geojson", false, "Emit Who's On First records as GeoJSON records where the 'properties' element is replaced by a Standard Place Response (SPR) representation of the record.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Emit one or more Who's On First records.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s emit [options] path(N) path(N)\n", os.Args[0])
		fs.PrintDefaults()
	}

	return fs
}
