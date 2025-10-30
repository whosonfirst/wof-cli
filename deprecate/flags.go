package supersede

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
)

var superseded_by_id int64
var superseded_by_reader_uri string
var superseded_by_writer_uri string

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("deprecated")

	fs.Int64Var(&superseded_by_id, "superseded-by-id", -1, "The ID to supersede each record with. If `-1` then this flag will be ignored.")

	fs.StringVar(&superseded_by_reader_uri, "superseded-by-reader-uri", "null://", "A valid whosonfirst/go-reader URI used to load records that are doing the superseding. Required if -superseding-id is not `-1`.")
	fs.StringVar(&superseded_by_writer_uri, "superseded-by-writer-uri", "null://", "A valid whosonfirst/go-writer URI used to update records that are doing the superseding. Required if -superseding-id is not `-1`.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Deprecate one or more Who's On First records.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options] path(N) path(N)\n", os.Args[0])
		fs.PrintDefaults()
	}

	return fs
}
