# wof

## wof remove-property

Ensure typed values for one or more Who's On First records.

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

### Examples

For example:

```
$> ./bin/wof remove-property \
	-property properties.media:properties.palette \
	-iterator-uri 'repo://' \
	-writer-uri repo:///usr/local/data/sfomuseum-data-media-collection \
	/usr/local/data/sfomuseum-data-media-collection	
```