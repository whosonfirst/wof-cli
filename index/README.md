# wof

## index

Index one or more Who's On First style data sources.

```
$> ./bin/wof index -h
Usage: wof index [TARGET] [OPTIONS]
Valid commands are:
* sql
```

_Support for indexing data sources in to an OpenSearch database is available in a separate [whosonfirst/go-whosonfirst-opensearch](https://github.com/whosonfirst/go-whosonfirst-opensearch) package which will, eventually, also be available from this tool._

### sql

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

### Database support

As of this writing the `index sql` tool has support for three databases: SQLite, MySQL and Postgres. Database support is NOT enabled by default and needs to be specified when you build the `wof` tool. These defaults may change in future releases but for the time being that's how it works.

| Database | Build tag | Notes |
| --- | --- | --- |
| SQLite | `sqlite3` | If you are indexing a database with either the "search" or "spelunker" tables you will also need to enable support for the FTS5 extension with the `fts5` build tag. |
| MySQL | `mysql` | |
| Postgres | `postgres` | |

For example:

```
$> cd wof-cli
$> make cli TAGS=sqlite3,fts5,postgres
```

### Database tables

Database tables and their schemas are defined in the [whosonfirst/go-whosonfirst-database/sql/tables](https://github.com/whosonfirst/go-whosonfirst-database/tree/main/sql/tables) package.

## Examples

Create a SQLite database with both the "spelunker" and "spatial" tables derived from three separate [sfomuseum-data](https://github.com/sfomuseum-data) repositories.

```
$> ./bin/wof index sql \
	-database-uri 'sql://sqlite3?dsn=test2.db' \
	-spatial-tables \
	-spelunker-tables \
	/usr/local/data/sfomuseum-data-publicart \
	/usr/local/data/sfomuseum-data-architecture \
	/usr/local/data/sfomuseum-data-whosonfirst
	
2025/11/05 09:44:15 INFO Iterator stats elapsed=24.544807417s seen=4501 allocated="19 MB" "total allocated"="18 GB" sys="513 MB" numgc=694
```

_Note that indexing databases with the "spelunker" or "search" tables can take a long time for large data sources (for example `whosonfirst-data-admin-us`) because of the time it takes to index the search table._