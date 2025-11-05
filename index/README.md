# wof

## index

### sql

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

### Examples

```
$> ./bin/wof index sql \
	-spatial-tables \
	-processes 100 \
	-database-uri 'sql://sqlite3?dsn=test.db' \
	/usr/local/data/whosonfirst-data-admin-us

2025/11/04 15:59:15 INFO Iterator stats elapsed=1m0.000103458s seen=87054 allocated="3.1 MB" "total allocated"="14 GB" sys="51 MB" numgc=4166
2025/11/04 16:00:15 INFO Iterator stats elapsed=2m0.000590291s seen=213928 allocated="4.9 MB" "total allocated"="27 GB" sys="135 MB" numgc=8454
2025/11/04 16:01:15 INFO Iterator stats elapsed=3m0.000484375s seen=321211 allocated="4.3 MB" "total allocated"="45 GB" sys="135 MB" numgc=17088
2025/11/04 16:02:15 INFO Iterator stats elapsed=4m0.0007335s seen=421566 allocated="5.2 MB" "total allocated"="61 GB" sys="135 MB" numgc=24467
2025/11/04 16:02:43 INFO Iterator stats elapsed=4m28.389431541s seen=448570 allocated="6.0 MB" "total allocated"="65 GB" sys="135 MB" numgc=26423
```