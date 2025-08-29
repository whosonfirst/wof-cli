package repos

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
)

var org string

var prefix multi.MultiString
var exclude multi.MultiString

var updated_since string
var forked bool
var not_forked bool

var token string

var exclude_archived bool
var ensure_commits bool

var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("repos")

	fs.StringVar(&org, "org", "whosonfirst-data", "The name of the organization to clone repositories from")

	fs.Var(&prefix, "prefix", "Limit repositories to only those with this prefix")
	fs.Var(&exclude, "exclude", "Exclude repositories with this prefix")

	fs.StringVar(&updated_since, "updated-since", "", "A valid Unix timestamp or an ISO8601 duration string (months are currently not supported)")
	fs.BoolVar(&forked, "forked", false, "Only include repositories that have been forked")
	fs.BoolVar(&not_forked, "not-forked", false, "Only include repositories that have not been forked")
	fs.StringVar(&token, "token", "", "A valid GitHub API access token")

	fs.BoolVar(&exclude_archived, "exclude-archived", false, "Exclude repos that have been archived.")

	fs.BoolVar(&ensure_commits, "ensure-commits", false, "Ensure that 1 or more files have been updated in the last commit")

	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "List repositories for the Who's On First (or other GitHub) organization.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options]\n", os.Args[0])
		fs.PrintDefaults()
	}

	return fs
}
