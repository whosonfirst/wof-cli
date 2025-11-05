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
$> ./bin/wof spelunker \
	-spelunker-uri 'sql://sqlite3?dsn=file:/usr/local/whosonfirst/wof-cli/work/test.db' \
	-protomaps-api-key {PROTOMAPS_API_KEY}

2025/11/04 18:13:50 INFO Listening for requests address=http://localhost:8080
```
