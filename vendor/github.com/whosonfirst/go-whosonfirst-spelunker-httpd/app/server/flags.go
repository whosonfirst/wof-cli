package server

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var server_uri string
var spelunker_uri string
var authenticator_uri string
var protomaps_api_key string

var root_url string

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("spelunker")

	fs.StringVar(&server_uri, "server-uri", "http://localhost:8080", "A valid `aaronland/go-http/v3/server.Server URI.")
	fs.StringVar(&spelunker_uri, "spelunker-uri", "null://", "A URI in the form of 'sql://{DATABASE_SQL_ENGINE}?dsn={DATABASE_SQL_DSN}' referencing the underlying Spelunker database. For example: sql://sqlite3?dsn=spelunker.db")
	fs.StringVar(&authenticator_uri, "authenticator-uri", "null://", "A valid aaronland/go-http/v3/auth.Authenticator URI. This is future-facing work and can be ignored for now.")
	fs.StringVar(&protomaps_api_key, "protomaps-api-key", "", "A valid Protomaps API key for displaying maps.")
	fs.StringVar(&root_url, "root-url", "", "The root URL for all public-facing URLs and links. If empty then the value of the -server-uri flag will be used.")

	return fs
}
