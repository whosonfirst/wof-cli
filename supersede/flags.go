package supersede

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
)

var superseding_reader_uri string
var superseding_writer_uri string
var parent_reader_uri string

var superseding_id int64
var parent_id int64

var is_deprecated bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("supersede")
	fs.Int64Var(&parent_id, "parent-id", -1, "The parent ID to assign to the new record.")
	fs.Int64Var(&superseding_id, "superseding-id", -1, "The ID to supersede each record with. If `-1` then each record will be cloned and the new ID of the clone will be used as the superseding_id.")

	fs.StringVar(&superseding_reader_uri, "superseding-reader-uri", "null://", "A valid whosonfirst/go-reader URI used to load records that are doing the superseding. Required if -superseding-id is not `-1`.")
	fs.StringVar(&superseding_writer_uri, "superseding-writer-uri", "null://", "A valid whosonfirst/go-writer URI used to update records that are doing the superseding. Required if -superseding-id is not `-1`.")
	fs.StringVar(&parent_reader_uri, "parent-reader-uri", "null://", "A valid whosonfirst/go-reader URI used to load parent records if -superseded-id is `-1`. Required if -parent-id is not `-1`.")

	fs.BoolVar(&is_deprecated, "deprecated", false, "Each record being superseded should also be marked as deprecated.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Supersede one or more Who's On First record and update the records superseding them.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options] path(N) path(N)\n", os.Args[0])
		fs.PrintDefaults()
	}

	return fs
}
