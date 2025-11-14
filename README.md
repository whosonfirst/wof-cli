# wof-cli

Command-line tool for common Who's On First operations.

## Tools

```
$> make cli
go build -mod vendor -ldflags="-s -w" -o bin/wof cmd/wof/main.go
```

Features and functionality are enable through the use of (Go language) [build tags](https://pkg.go.dev/cmd/go#hdr-Build_constraints). By default all build tags are enabled. Please consult the [Build tags](#build-tags) section below for detailed descriptions.

### wof

```
$> ./bin/wof -h
Usage: wof [CMD] [OPTIONS]
Valid commands are:
* centroid
* deprecate
* edit
* emit
* ensure-property
* export
* fmt
* geometry
* index
* open
* pip
* property
* rebuild-hierarchy
* remove-property
* repos
* show
* spelunker
* supersede
* uri
* validate
```

_Important: The inputs and outputs for the `wof` tool have not been finalized yet, notably about how files are read and written if updated. You should expect change in the short-term._

#### wof centroid

```
$> ./bin/wof centroid -h
Emit centroids data (source, latitude, longitude) for one or more Who's On First records as CSV-encoded properties.
Usage:
	 ./bin/wof emit [options] path(N) path(N)
```

Examples can be found in [centroid/README.md](centroid/README.md#examples)

#### wof deprecate

```
$> ./bin/wof deprecate -h
Deprecate one or more Who's On First records.
Usage:
	 ./bin/wof [options] path(N) path(N)
  -superseded-by-id -1
    	The ID to supersede each record with. If -1 then this flag will be ignored. (default -1)
  -superseded-by-reader-uri -1
    	A valid whosonfirst/go-reader URI used to load records that are doing the superseding. Required if -superseding-id is not -1. (default "null://")
  -superseded-by-writer-uri -1
    	A valid whosonfirst/go-writer URI used to update records that are doing the superseding. Required if -superseding-id is not -1. (default "null://")
```

Examples can be found in [deprecate/README.md](deprecate/README.md#examples)

#### wof edit

```
$> ./bin/wof edit -h
Launch a local web server running on a random port hosting a web application for editing one or more Who's On First records.
Usage:
	 ./bin/wof path(N) path(N)
  -ensure-relative-path
    	Boolean flag signaling that each URI should be expanded to its fully-quality WOF-style relative path. This flag is only processed if the -reader-uri flag is not-empty.
  -gh-access-token-uri string
    	A valid GitHub API access token. This is only necessary if -writer-uri is "wof-pr://".
  -map-provider string
    	Valid options are: leaflet, protomaps (default "leaflet")
  -map-tile-uri string
    	A valid Leaflet tile layer URI. See documentation for special-case (interpolated tile) URIs. (default "https://tile.openstreetmap.org/{z}/{x}/{y}.png")
  -protomaps-max-data-zoom int
    	The maximum zoom (tile) level for data in a PMTiles database. Necessary for "over-zooming".
  -protomaps-theme string
    	A valid Protomaps theme label. (default "light")
  -reader-uri string
    	An optional whosonfirst/go-reader/v2.Reader URI used to read WOF records from alternate sources. If defined then the -writer-uri flag must also be populated.
  -verbose
    	Enable verbose (debug) logging.
  -writer-uri string
    	An optional whosonfirst/go-writer/v3.Writer URI used to write records to alternate sources. If defined then the -reader-uri flag must also be populated.
```

Examples can be found in [edit/README.md](edit/README.md#examples)

#### wof emit

```
$> ./bin/wof emit -h
Emit one or more Who's On First records.
Usage:
	 ./bin/wof emit [options] path(N) path(N)
  -as-spr
    	Emit Who's On First records formatted as Standard Place Response (SPR) records. This flag is DEPRECATED. Please use '-format spr' instead.
  -as-spr-geojson
    	Emit Who's On First records as GeoJSON records where the 'properties' element is replaced by a Standard Place Response (SPR) representation of the record. This flag is DEPRECATED. Please use '-format geojson' instead.
  -csv-append-property value
    	Zero or more additional properties to append to each CSV row. Properties should be in the format of {COLUMN_NAME}={PATH}. This flag is only honoured if the -format flag has a value of "csv".
  -forgiving
    	Do not stop processing when errors are encountered.
  -format string
    	Valid options are: csv, spr, spr-geojson or [none]. If none then the raw GeoJSON for each matching record will be emitted.
  -include-alt-geoms
    	Emit alternate geometry records. (default true)
  -iterator-uri string
    	A valid whosonfirst/go-whosonfirst-iterate/v2/emitter URI. If URI is "-" then this flag will be assigned a value of "file://" whose input will be the expanded URIs derived from additional arguments. Available options are: cwd://, directory://, featurecollection://, file://, filelist://, geojsonl://, null://, repo:// (default "repo://")
  -query value
    	One or more {PATH}={REGEXP} parameters for filtering records.
  -query-mode string
    	Specify how query filtering should be evaluated. Valid modes are: ALL, ANY (default "ALL")
  -writer-uri string
    	A valid whosonfirst/go-writer.Writer URI. Available options are: cwd://, featurecollection://, fs://, geoparquet://, io://, jsonl://, null://, repo://, sqlite://, stdout:// (default "jsonl://?writer=stdout://")
```

Examples can be found in [emit/README.md](emit/README.md#examples)

#### wof ensure-property

Ensure typed values for one or more Who's On First records.

```
$> ./bin/wof ensure-property -h
Ensure typed values for one or more Who's On First records.
Usage:
	./bin/wof [options] path(N) path(N)
Valid options are:
  -bool-property-from string
    	A valid tidwall/gjson path used to erive the value of the -bool-property from a property in the same document.
  -boolean-property value
    	One or more {KEY}={VALUE} flags where {KEY} is a valid tidwall/gjson path and {VALUE} is a boolean value.
  -exporter-uri string
    	A valid whosonfirst/go-whosonfirst-export/v3.Exporter URI. (default "whosonfirst://")
  -float-property value
    	One or more {KEY}={VALUE} flags where {KEY} is a valid tidwall/gjson path and {VALUE} is a float(64) value.
  -float-property-from string
    	A valid tidwall/gjson path used to erive the value of the -float-property from a property in the same document.
  -geometry-property value
    	A {KEY}={VALUE} flag indicating the source of the geometry data to assign. Valid options are: wkt={VALID_WKT_GEOMETRY}, geojson={VALID_GEOJSON_GEOMETRY}, file={PATH_TO_GEOJSON_FEATURE}.
  -if-missing
    	Only assign property value if the property key is not set.
  -int-property value
    	One or more {KEY}={VALUE} flag where {KEY} is a valid tidwall/gjson path and {VALUE} is a int(64) value.
  -int-property-from string
    	A valid tidwall/gjson path used to erive the value of the -int-property from a property in the same document.
  -iterator-uri string
    	A valid whosonfirst/go-whosonfirst-iterate/v3.Iterator URI. (default "repo://")
  -string-property value
    	One or more {KEY}={VALUE} flag where {KEY} is a valid tidwall/gjson path and {VALUE} is a string value.
  -string-property-from string
    	A valid tidwall/gjson path used to erive the value of the -string-property from a property in the same document.
  -verbose
    	Enable verbose (debug) logging
  -writer-uri string
    	A valid whosonfirst/go-writer/v3.Writer URI. (default "null://")
```

Examples can be found in [property/ensure/README.md](property/ensure/README.md#examples)

#### wof export

```
> ./bin/wof export -h
"Export-ify" one or more Who's On First records.
Usage:
	 ./bin/wof path(N) path(N)
  -stdout
    	Boolean flag signaling that updated records should be written to STDOUT. If false input files will be overwritten.
```

#### wof fmt

```
$> ./bin/wof fmt -h
Format one or more GeoJSON files according to Who's On First conventions.
Usage:
	 ./bin/wof path(N) path(N)
  -stdout
    	Boolean flag signaling that updated records should be written to STDOUT. If false input files will be overwritten. (default true)
```

Examples can be found in [format/README.md](format/README.md#examples)

#### wof geometry

```
$> ./bin/wof geometry -h
Display or update the geometry for one or more Who's On First records.
Usage:
	 ./bin/wof geometry [options] path(N) path(N)
  -action string
    	Valid options are: show, update. (default "show")
  -format string
    	The format of the source geometry (if -action update). Valid options are: geojson. (default "geojson")
  -source string
    	The data source from which to derive a geometry (if -action update). Valid options are: the path to a local file.
  -stdout
    	Boolean flag signaling that updated records should be written to STDOUT. If false input files will be overwritten.
```

Examples can be found in [geometry/README.md](geometry/README.md#examples)

#### wof index

Index one or more Who's On First style data sources.

```
$> ./bin/wof index -h
Usage: wof index [TARGET] [OPTIONS]
Valid commands are:
* sql
```

##### wof index sql

Index one or more Who's On First style data sources in to a [database/sql](https://pkg.go.dev/database/sql) compatible database.

```
$> ./bin/wof index sql -h
  -all
    	Index all tables (except the 'search' and 'geometries' tables which you need to specify explicitly)
  -ancestors
    	Index the 'ancestors' tables
  -concordances
    	Index the 'concordances' tables
  -database-uri string
    	A URI in the form of 'sql://{DATABASE_SQL_ENGINE}?dsn={DATABASE_SQL_DSN}'. For example: sql://sqlite3?dsn=test.db
  -geojson
    	Index the 'geojson' table
  -geometries
    	Index the 'geometries' table (requires that libspatialite already be installed)
  -index-alt value
    	Zero or more table names where alt geometry files should be indexed.
  -index-relations
    	Index the records related to a feature, specifically wof:belongsto, wof:depicts and wof:involves. Alt files for relations are not indexed at this time.
  -index-relations-reader-uri string
    	A valid go-reader.Reader URI from which to read data for a relations candidate.
  -iterator-uri string
    	A valid whosonfirst/go-whosonfirst-iterate/v3.Iterator URI. Supported iterator URI schemes are: cwd://,directory://,featurecollection://,file://,filelist://,geojsonl://,git://,githubapi://,githuborg://,null://,reader://,repo:// (default "repo://")
  -names
    	Index the 'names' table
  -optimize
    	Attempt to optimize the database before closing connection (default true)
  -processes int
    	The number of concurrent processes to index data with (default 20)
  -properties
    	Index the 'properties' table
  -rtree
    	Index the 'rtree' table
  -search
    	Index the 'search' table (using SQLite FTS4 full-text indexer)
  -spatial-tables
    	If true then index the necessary tables for use with the whosonfirst/go-whosonfirst-spatial-sqlite package.
  -spelunker
    	Index the 'spelunker' table
  -spelunker-tables
    	If true then index the necessary tables for use with the whosonfirst/go-whosonfirst-spelunker packages
  -spr
    	Index the 'spr' table
  -strict-alt-files
    	Be strict when indexing alt geometries (default true)
  -supersedes
    	Index the 'supersedes' table
  -verbose
    	Enable verbose (debug) logging
```

Examples can be found in [index/README.md](index/README.md#examples)

#### wof open

Open one or more Who's On First documents in a custom editor.

```
$> ./bin/wof open -h
Open one or more Who's On First documents in a custom editor.
Usage:
	 ./bin/wof path(N) path(N)
  -editor string
    	The editor to use for opening Who's On First records. If empty the value of the EDITOR environment variable will be used.
```

#### wof pip

Perform point-in-polygon and wof:hierarchy update operations on one or more Who's On First records.

```
$> ./bin/wof pip -h
Perform point-in-polygon and wof:hierarchy update operations on one or more Who's On First records.
Usage:
	 ./bin/wof path(N) path(N)
  -export
    	"Export-ify" each record after point-in-polygon operations are complete.
  -mapshaper-client-uri string
    	Optional URI to a sfomuseum/go-sfomuseum-mapshaper server instance used to derive point-in-polygon centroids. If absent then the centroid used to perform point-in-polygon operations will be determined using internal heuristics.
  -placetype string
    	Assign this value as the "wof:placetype" property before performing point-in-polygon operations.
  -spatial-database-uri string
    	A valid whosonfirst/go-whosonfirst-spatial/database URI. By default 'pmtiles://' and 'sqlite://' spatial database URIs are supported.
  -stdout
    	Boolean flag signaling that updated records should be written to STDOUT. (default true)
```

Examples can be found in [pip/README.md](pip/README.md#examples)

#### wof property

Print one or more properties for one or more Who's On First IDs.

```
$> wof property -h
Print one or more properties for one or more Who's On First IDs.
Usage:
	 wof path(N) path(N)
  -format string
    	Valid options are: csv. If empty then properties will printed as a new-line separated list.
  -path value
    	One or more valid tidwall/gjson paths to extract from each document
  -prefix string
    	If not empty this prefix will be appended (and separated by a ".") to each -path argument
```

Examples can be found in [property/README.md](property/README.md#examples)

#### wof rebuild-hierarchy

Rebuild the wof:hierarchy property one or more Who's On First records.

```
$> ./bin/wof rebuild-hierarchy -h
Rebuild the wof:hierarchy property one or more Who's On First records.
Usage:
	./bin/wof [options] path(N) path(N)
Valid options are:
  -exporter-uri string
    	A valid whosonfirst/go-whosonfirst-export/v3.Exporter URI. (default "whosonfirst://")
  -iterator-uri string
    	A valid whosonfirst/go-whosonfirst-iterate/v3.Iterator URI. (default "repo://")
  -parent-id value
    	One or more explicit parent IDs to use for deriving hierarchies. The default is to use each record's wof:parent_id value.
  -parent-reader-uri string
    	A valid whosonfirst/go-reader/v2.Reader URI. (default "https://data.whosonfirst.org")
  -verbose
    	Enable verbose (debug) logging
  -writer-uri string
    	A valid whosonfirst/go-writer/v3.Writer URI. (default "null://")
```	

#### wof remove-property

Remove properties from one or more Who's On First records.

```
$> ./bin/wof remove-property -h
Remove properties from one or more Who's On First records.
Usage:
	./bin/wof [options] path(N) path(N)
Valid options are:
  -exporter-uri string
    	A valid whosonfirst/go-whosonfirst-export/v3.Exporter URI. (default "whosonfirst://")
  -iterator-uri string
    	A valid whosonfirst/go-whosonfirst-iterate/v3.Iterator URI. (default "repo://")
  -property value
    	One or more (fully-qualified) paths for a property to remove.
  -verbose
    	Enable verbose (debug) logging
  -writer-uri string
    	A valid whosonfirst/go-writer/v3.Writer URI. (default "null://")
```

Examples can be found in [property/remove/README.md](property/remove/README.md#examples)

#### wof repos

List repositories for the Who's On First (or other GitHub) organization.

```
$> ./bin/wof repos -h
List repositories for the Who's On First (or other GitHub) organization.
Usage:
	 ./bin/wof [options]
  -ensure-commits
    	Ensure that 1 or more files have been updated in the last commit
  -exclude value
    	Exclude repositories with this prefix
  -exclude-archived
    	Exclude repos that have been archived.
  -forked
    	Only include repositories that have been forked
  -not-forked
    	Only include repositories that have not been forked
  -org string
    	The name of the organization to clone repositories from (default "whosonfirst-data")
  -prefix value
    	Limit repositories to only those with this prefix
  -token string
    	A valid GitHub API access token
  -updated-since string
    	A valid Unix timestamp or an ISO8601 duration string (months are currently not supported)
  -verbose
    	Enable verbose (debug) logging
```

Examples can be found in [repos/README.md](repos/README.md#examples)

#### wof show

Command-line tool for serving GeoJSON features from an on-demand web server.

```
$> ./bin/wof show -h
Command-line tool for serving GeoJSON features from an on-demand web server.
Usage:
	 ./bin/wof path(N) path(N)
Valid options are:
  -label value
    	Zero or more (GeoJSON Feature) properties to use to construct a label for a feature's popup menu when it is clicked on.
  -map-provider string
    	Valid options are: leaflet, protomaps (default "leaflet")
  -map-tile-uri string
    	A valid Leaflet tile layer URI. See documentation for special-case (interpolated tile) URIs. (default "https://tile.openstreetmap.org/{z}/{x}/{y}.png")
  -point-style string
    	A custom Leaflet style definition for point geometries. This may either be a JSON-encoded string or a path on disk.
  -port int
    	The port number to listen for requests on (on localhost). If 0 then a random port number will be chosen.
  -protomaps-theme string
    	A valid Protomaps theme label. (default "white")
  -style string
    	A custom Leaflet style definition for geometries. This may either be a JSON-encoded string or a path on disk.

If the only path as input is "-" then data will be read from STDIN.
```

Examples can be found in [show/README.md](show/README.md#examples)

#### wof spelunker

Launch a local [Who's On First Spelunker website](https://spelunker.whosonfirst.org) for browsing custom Who's On First style data.

```
$> ./bin/wof spelunker -h
  -authenticator-uri string
    	A valid aaronland/go-http/v3/auth.Authenticator URI. This is future-facing work and can be ignored for now. (default "null://")
  -protomaps-api-key string
    	A valid Protomaps API key for displaying maps.
  -root-url string
    	The root URL for all public-facing URLs and links. If empty then the value of the -server-uri flag will be used.
  -server-uri string
    	A valid `aaronland/go-http/v3/server.Server URI. (default "http://localhost:8080")
  -spelunker-uri string
    	A URI in the form of 'sql://{DATABASE_SQL_ENGINE}?dsn={DATABASE_SQL_DSN}' referencing the underlying Spelunker database. For example: sql://sqlite3?dsn=spelunker.db (default "null://")
```

Examples can be found in [spelunker/README.md](spelunker/README.md#examples)

#### wof supersede

Supersede one or more Who's On First record and update the records superseding them.

```
$> ./bin/wof supersede -h
Supersede one or more Who's On First record and update the records superseding them.
Usage:
	 ./bin/wof [options] path(N) path(N)
  -deprecated
    	Each record being superseded should also be marked as deprecated.
  -parent-id int
    	The parent ID to assign to the new record. (default -1)
  -parent-reader-uri string
    	A valid whosonfirst/go-reader URI used to load parent records if -superseded-id is -1. Required if -parent-id is not `-1`. (default "null://")
  -superseding-id int
    	The ID to supersede each record with. If -1 then each record will be cloned and the new ID of the clone will be used as the superseding_id. (default -1)
  -superseding-reader-uri string
    	A valid whosonfirst/go-reader URI used to load records that are doing the superseding. Required if -superseding-id is not -1. (default "null://")
  -superseding-writer-uri string
    	A valid whosonfirst/go-writer URI used to update records that are doing the superseding. Required if -superseding-id is not -1. (default "null://")
```

Examples can be found in [supersede/README.md](supersede/README.md#examples)

#### wof uri

Print the nested URI for one or more Who's On First IDs.

```
$> ./bin/wof uri -h
Print the nested URI for one or more Who's On First IDs.
Usage:
	 ./bin/wof path(N) path(N)
  -prefix string
    	An optional prefix to append to the final URI.
```

Examples can be found in [uri/README.md](uri/README.md#examples)

#### wof validate

Validate one or more Who's On First documents.

```
$> ./bin/wof validate -h
Validate one or more Who's On First documents.
Usage:
	 ./bin/wof path(N) path(N)
```

Examples can be found in [validate/README.md](validate/README.md#examples)

## Paths and URIs

The default behaviour for the `wof` command is to assume that all the URIs it is passed are paths on the local computer.

That being said all the tool also "expand" each path using the [uris.ExpandURIsWithCallback](uris/uris.go) method which allows for more sophisticated behaviour in the future.

### Expansions

Currently there is only two supported "expansions":

1. Treating bare numbers as Who's On First IDs, resolving them to their relative path and looking for that file in a parent "data" directory inside the current working directory. Basically it's a shortcut for resolving a Who's On First record to its GeoJSON representation inside a `whosonfirst-data` repository.
2. URIs prefixed with `repo://` will be treated as a [whosonfirst/go-whosonfirst-iterate/v2 "repo" emitter URI](https://github.com/whosonfirst/go-whosonfirst-iterate?tab=readme-ov-file#repo) and all the files in that repository will be processed.

Eventually there may be other "expansions" most notably support for the Go `./...` syntax to process all the Who's On First records in the current working directory.

## Advanced

### Build tags

Features and functionality are enable through the use of (Go language) build tags. By default the proverbial "kitchen sink" (meaning all features and functionality) is enabled but you can compile your own more limited copy of the `wof` through the use of `no_{FEATURE}` build tags.

For example to build the `wof` command line tool without suport for the [whosonfirst/go-writer-geoparquet](https://github.com/whosonfirst/go-writer-geoparquet) package (which would shave about 20MB off the size of the final binary) you would do this:

```
$> go build -mod vendor -ldflags="-s -w" -tags no_writer_geoparquet -o bin/wof cmd/wof/main.go
```

#### Available build tags

##### no_centroid

Disable the `wof centroid` command.

##### no_emit

Disable the `wof emit` command.

##### no_ensure_property

Disable the `wof ensure-property` command.

##### no_export

Disable the `wof export` command.

##### no_format

Disable the `wof format` command.

##### no_geometry

Disable the `wof geometry` command.

##### no_iterator_git

Disable import of the [whosonfirst/go-whosonfirst-iterate-git/v3](https://github.com/whosonfirst/go-whosonfirst-iterate-git) package.

This affects the `wof emit` command.

##### no_iterator_org

Disable import of the [whosonfirst/go-whosonfirst-iterate-organization/v3](https://github.com/whosonfirst/go-whosonfirst-iterate-organization) package.

This affects the `wof emit` command.

##### no_open

Disable the `wof open` command.

##### no_pip

Disable the `wof pip` command.

##### no_property

Disable the `wof property` command.

##### no_remove_property

Disable the `wof remove-property` command.

##### no_repos

Disable the `wof repos` command.

##### no_show

Disable the `wof show` command.

##### no_spatial_pmtiles

Disable import of the [whosonfirst/go-whosonfirst-spatial-pmtiles](https://github.com/whosonfirst/go-whosonfirst-spatial-pmtiles) package.

This affects the `wof pip` command.

##### no_spatial_sqlite

Disable import of the [whosonfirst/go-whosonfirst-spatial-sqlite](https://github.com/whosonfirst/go-whosonfirst-spatial-sqlite) package.

This affects the `wof pip` command.

##### no_supersede

Disable the `wof supersede` command.

##### no_uri

Disable the `wof uri` command.

##### no_validate

Disable the `wof validate` command.

##### no_writer_geoparquet

Disable import of the [whosonfirst/go-writer-geoparquet](https://github.com/whosonfirst/go-writer-geoparquet) package.

This affects the `wof emit` command.

## See also

* https://github.com/whosonfirst/go-whosonfirst-export
* https://github.com/whosonfirst/go-whosonfirst-format
* https://github.com/whosonfirst/go-whosonfirst-github
* https://github.com/whosonfirst/go-whosonfirst-validate
* https://github.com/whosonfirst/go-whosonfirst-spatial
* https://github.com/whosonfirst/go-whosonfirst-spatial-pmtiles
* https://github.com/whosonfirst/go-whosonfirst-spatial-sqlite
* https://github.com/whosonfirst/go-whosonfirst-iterate
* https://github.com/sfomuseum/go-sfomuseum-mapshaper
