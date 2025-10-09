# wof

## wof rebuild-hierarchy

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

### Examples

For example:

```
$> /bin/wof rebuild-hierarchy \
	-verbose \
	-iterator-uri file:/// \	
	-parent-reader-uri https://static.sfomuseum.org \
	-parent-id 102527513 \
	-writer-uri repo:///usr/local/data/sfomuseum-data-publicart \
	data/137/745/518/3/1377455183.geojson \
	data/137/745/522/5/1377455225.geojson \
	data/186/348/188/5/1863481885.geojson \
	data/137/745/522/5/1377455225.geojson
	
2025/10/09 14:48:48 DEBUG Verbose logging enabled
2025/10/09 14:48:48 DEBUG Do iter uri=data/137/745/522/5/1377455225.geojson attempt=0 "max attempts"=1 counter=0
2025/10/09 14:48:48 DEBUG Iterate target uri=data/137/745/522/5/1377455225.geojson "target uri"=data/137/745/522/5/1377455225.geojson counter=0 "local counter"=0
2025/10/09 14:48:48 DEBUG Do iter uri=data/137/745/518/3/1377455183.geojson attempt=0 "max attempts"=1 counter=0
2025/10/09 14:48:48 DEBUG Iterate target uri=data/137/745/518/3/1377455183.geojson "target uri"=data/137/745/518/3/1377455183.geojson counter=0 "local counter"=0
2025/10/09 14:48:48 DEBUG Iteration successful uri=data/137/745/522/5/1377455225.geojson attempt=1 "max attempts"=1 counter=1
2025/10/09 14:48:48 DEBUG Do iter uri=data/186/348/188/5/1863481885.geojson attempt=0 "max attempts"=1 counter=0
2025/10/09 14:48:48 DEBUG Iterate target uri=data/186/348/188/5/1863481885.geojson "target uri"=data/186/348/188/5/1863481885.geojson counter=0 "local counter"=0
2025/10/09 14:48:48 DEBUG Do iter uri=data/137/745/522/5/1377455225.geojson attempt=0 "max attempts"=1 counter=0
2025/10/09 14:48:48 DEBUG Iterate target uri=data/137/745/522/5/1377455225.geojson "target uri"=data/137/745/522/5/1377455225.geojson counter=0 "local counter"=0
2025/10/09 14:48:48 INFO Updated record path=data/137/745/522/5/1377455225.geojson
2025/10/09 14:48:48 DEBUG Iteration successful uri=data/137/745/518/3/1377455183.geojson attempt=1 "max attempts"=1 counter=1
2025/10/09 14:48:48 INFO Updated record path=data/137/745/518/3/1377455183.geojson
2025/10/09 14:48:48 DEBUG Iteration successful uri=data/137/745/522/5/1377455225.geojson attempt=1 "max attempts"=1 counter=1
2025/10/09 14:48:48 DEBUG Run garbage collector uri=data/137/745/522/5/1377455225.geojson
2025/10/09 14:48:48 DEBUG Time to iterate uri uri=data/137/745/522/5/1377455225.geojson time=682.656542ms
2025/10/09 14:48:49 INFO Updated record path=data/137/745/522/5/1377455225.geojson
2025/10/09 14:48:49 DEBUG Iteration successful uri=data/186/348/188/5/1863481885.geojson attempt=1 "max attempts"=1 counter=1
2025/10/09 14:48:49 DEBUG Run garbage collector uri=data/137/745/518/3/1377455183.geojson
2025/10/09 14:48:49 DEBUG Time to iterate uri uri=data/137/745/518/3/1377455183.geojson time=1.032810958s
2025/10/09 14:48:49 INFO Updated record path=data/186/348/188/5/1863481885.geojson
2025/10/09 14:48:49 DEBUG Time to process paths count=4 time=1.355403083s
2025/10/09 14:48:49 DEBUG Run garbage collector uri=data/137/745/522/5/1377455225.geojson
```