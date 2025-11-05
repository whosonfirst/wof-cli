# go-whosonfirst-spelunker

Go package implementing a common interface for Who's On First "spelunker"-ing.

## Documentation

Documentation is incomplete at this time.

## Important

This is work in progress and you should expect things to change, break or simply not work yet.

## Motivation

This is a refactoring of both the [whosonfirst/whosonfirst-www-spelunker](github.com/whosonfirst/whosonfirst-www-spelunker) and [whosonfirst/go-whosonfirst-browser](github.com/whosonfirst/go-whosonfirst-browser) packages.

Specifically, the former (`whosonfirst-www-spelunker`) is written in Python and has a sufficiently complex set of requirements that spinning up a new instance is difficult. By rewriting the spelunker tool in Go the hope is to eliminate or at least minimize these external requirements and to make it easier to deploy the spelunker to "serverless" environments like AWS Lambda or Function URLs. The latter (`go-whosonfirst-browser`) has developed a sufficiently large and complex code base that starting from scratch and simply copying, and adapting, existing functionality seemed easier than trying to refactor everything.

## A note about versioning

Currently this package is unversioned reflecting the fact that it is still in flux. The rate of change is slowing down and will eventually be assigned version numbers less than 1.x for as long as it takes to produce the initial "minimal viable (and working)" Spelunker implementations. These versions (0.x.y) should not be considered to be backwards compatible with each other and are expected to change as the first stable interface is settled, specifically if and whether it will contain spatial functions.

Once a decision has been reached on that matter and everything is proven to work this package (and all the related packages, discussed below) will be bumped up to a "version 2.x" release, skipping version 1.x altogether, reflecting the fact that the original Python version of the Spelunker is "version 1" and that this code base is meaningfully different.

After the "v2" release this package (and related packages) will follow the standard Go convention of incrementing version numbers if and when there are changes to the underlying Spelunker interface.

## Structure

There are three "classes" of `go-whosonfirst-spelunker` packages:

### go-whosonfirst-spelunker

That would be the package that you are looking at right now. It defines the [Spelunker](https://github.com/whosonfirst/go-whosonfirst-spelunker/blob/main/spelunker.go#L24-L82) interface defining the minimal methods required for a Spelunker application be it a command-line application, a web application or something else.

This package does not export any working implementations of the `Spelunker` interface. It simply defines the interface and other associated types.

### go-whosonfirst-spelunker-httpd

The [whosonfirst/go-whosonfirst-spelunker-httpd](github.com/whosonfirst/go-whosonfirst-spelunker-httpd) package provides libraries for implementing a web-based spelunker service. While it does define a working `cmd/server` tool demonstrating how those libraries can be used, like the `go-whosonfirst-spelunker` package it does not export any working implementations of the `Spelunker` interface. 

The idea is to separate the interaction details and the mechanics of a web application from the details of how data is stored or queried from any given database containing Who's On First records. 

The server itself can be run and will serve requests because its default database is the [NullSpelunker](https://github.com/whosonfirst/go-whosonfirst-spelunker/blob/main/spelunker_null.go) implementation but since that implementation simply returns "Not implemented" for every method in the Spelunker interface it probably won't be of much use.

### go-whosonfirst-spelunker-{DATABASE}

These are packages that implement the [Spelunker](https://github.com/whosonfirst/go-whosonfirst-spelunker/blob/main/spelunker.go#L24-L82) interface for a particular database engine. Current implementations are:

#### go-whosonfirst-spelunker-opensearch

The [whosonfirst/go-whosonfirst-spelunker-opensearch](github.com/whosonfirst/go-whosonfirst-spelunker-opensearch) package implements the `Spelunker` interface using an [OpenSearch](https://opensearch.org/) document store, for example data indexed by the [whosonfirst/go-whosonfirst-opensearch](https://github.com/whosonfirst/go-whosonfirst-opensearch) package.

It imports both the `go-whosonfirst-spelunker` and `go-whosonfirst-spelunker-httpd` and exports local instances of the web-based server (`httpd`).

_Set up and example(s) to be written..._

#### go-whosonfirst-spelunker-sql

The [whosonfirst/go-whosonfirst-spelunker-sqlite](github.com/whosonfirst/go-whosonfirst-spelunker-sql) package implements the `Spelunker` interface using a Go `database/sql` relational database source, for example SQLite databases produced by the [whosonfirst/go-whosonfirst-sqlite-features-index](https://github.com/whosonfirst/go-whosonfirst-sqlite-features-index) package.

It imports both the `go-whosonfirst-spelunker` and `go-whosonfirst-spelunker-httpd` and exports local instances of the "spelunker" command-line tool and web-based server. For example, to create a database for use by the SQLite implementation of the `Spelunker` interface:

```
$> cd /usr/local/whosonfirst/go-whosonfirst-sqlite-features-index
$> ./bin/wof-sqlite-index-features-mattn \
	-timings \
	-database-uri mattn:///usr/local/data/ca.db \
	-spelunker-tables \
	-index-alt geojson \
	/usr/local/data/whosonfirst-data-admin-ca
```

And then to use that database with a local (`go-whosonfirst-spelunker-sql`) instance of server code exported by the `go-whosonfirst-spelunker-httpd` package:

```
$> cd /usr/local/whosonfirst/go-whosonfirst-spelunker-sql
$> ./bin/server \
	-server-uri http://localhost:8080 \
	-spelunker-uri sql://sqlite3?dsn=file:/usr/local/data/ca.db
```

This is what the code for the server tool looks like (with error handling omitted for the sake of brevity):

```
package main

import (
        "context"
        "log/slog"

        _ "github.com/mattn/go-sqlite3"
        "github.com/whosonfirst/go-whosonfirst-spelunker-httpd/app/server"
        _ "github.com/whosonfirst/go-whosonfirst-spelunker-sql"
)

func main() {
        ctx := context.Background()
        logger := slog.Default()
        server.Run(ctx, logger)
}
```

_So far this package has only been tested with SQLite databases and probably contains some SQLite-specific syntax. The hope is that database engine specifics can be handled in conditionals in the `go-whosonfirst-spelunker-sql` package itself leaving consumers none the wiser._

#### go-whosonfirst-spelunker-sqlite

This package builds on the `whosonfirst/go-whosonfirst-spelunker-sql` and the `whosonfirst/go-whosonfirst-spelunker-httpd` packages but also imports @psanford 's [sqlite3vfs](https://github.com/psanford?tab=repositories&q=sqlite3vfs&type=&language=&sort=) packages to enable the use of SQLite databases hosted on remote servers.

It modifies the default `go-whosonfirst-spelunker-httpd` options and flags and then launches the spelunker server using the `RunWithOptions` method. For example (with error handling omitted for the sake of brevity):

```
import (
	"context"
	"fmt"
	"github.com/psanford/sqlite3vfs"
	"github.com/psanford/sqlite3vfshttp"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/whosonfirst/go-whosonfirst-spelunker-httpd/app/server"
	_ "github.com/whosonfirst/go-whosonfirst-spelunker-sql"
)

func main() {

	ctx := context.Background()
	logger := slog.Default()

	fs := server.DefaultFlagSet()

	opts, _ := server.RunOptionsFromFlagSet(ctx, fs)

	is_vfs, vfs_uri, _ := checkVFS(opts.SpelunkerURI)

	if is_vfs {
		opts.SpelunkerURI = vfs_uri
	}

	server.RunWithOptions(ctx, opts, logger)
}

func checkVFS(spelunker_uri string) (bool, string, error) {

	u, _ := url.Parse(spelunker_uri)

	q := u.Query()

	if !q.Has("vfs") {
		return false, spelunker_uri, nil
	}

	vfs_url := q.Get("vfs")

	vfs := sqlite3vfshttp.HttpVFS{
		URL: vfs_url,
		// Consult cmd/server/main.go for roundTripper; it has been
		// excluded here for the sake of brevity
		RoundTripper: &roundTripper{
			referer:   q.Get("vfs-referer"),
			userAgent: q.Get("vfs-user-agent"),
		},
	}

	sqlite3vfs.RegisterVFS("httpvfs", &vfs)

	dsn := "spelunker.db?vfs=httpvfs&mode=ro"
	q.Set("dsn", dsn)
	q.Del("vfs")

	u.RawQuery = q.Encode()
	return true, u.String(), nil
}
```

_Note: In practice this (querying a SQLite database over HTTP) doesn't really work in a Spelunker context. Specifically, it works for simple atomic queries but the moment the application starts to do multiple overlapping queries in the same session/context there are database locks and everything times out. Maybe I am doing something wrong? I would love to know what and how to fix it if that's the case since this is a super-compelling deployment strategy. Until then it should probably best be understood as a reference implementation only._

## See also

* https://github.com/whosonfirst/go-whosonfirst-spelunker-httpd
* https://github.com/whosonfirst/go-whosonfirst-spelunker-opensearch
* https://github.com/whosonfirst/go-whosonfirst-spelunker-sql
* https://github.com/whosonfirst/go-whosonfirst-spelunker-sqlite
