# wof

## spelunker

```
$> ./bin/wof spelunker -h
  -authenticator-uri string
    	... (default "null://")
  -protomaps-api-key string
    	...
  -root-url string
    	The root URL for all public-facing URLs and links. If empty then the value of the -server-uri flag will be used.
  -server-uri string
    	... (default "http://localhost:8080")
  -spelunker-uri string
    	... (default "null://")
```

## Examples

```
$> make cli TAGS=sqlite3,fts5
go build -mod vendor -ldflags="-s -w" -tags sqlite3,fts5 -o bin/wof cmd/wof/main.go
```

```
$> ./bin/wof index sql \
	-database-uri 'sql://sqlite3?dsn=work/test2.db' \
	-spatial-tables \
	-spelunker-tables \
	/usr/local/data/sfomuseum-data-publicart \
	/usr/local/data/sfomuseum-data-architecture \
	/usr/local/data/sfomuseum-data-whosonfirst
	
2025/11/05 15:17:45 INFO Iterator stats elapsed=29.556344625s seen=4501 allocated="23 MB" "total allocated"="18 GB" sys="572 MB" numgc=691
```

```
$> ./bin/wof spelunker \
	-spelunker-uri 'sql://sqlite3?dsn=file:/usr/local/whosonfirst/wof-cli/work/test2.db'
	-protomaps-api-key {PROTOMAPS_API_KEY}

2025/11/05 15:17:51 INFO Listening for requests address=http://localhost:8080
```
