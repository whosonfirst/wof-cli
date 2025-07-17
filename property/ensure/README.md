# wof

## wof ensure-property

Ensure typed values for one or more Who's On First records.

```
$> ./bin/wof ensure-property -h
Ensure typed values for one or more Who's On First records.
Usage:
	./bin/wof [options] path(N) path(N)
Valid options are:
  -boolean-property value
    	One or more {KEY}={VALUE} flags where {KEY} is a valid tidwall/gjson path and {VALUE} is a boolean value.
  -exporter-uri string
    	A valid whosonfirst/go-whosonfirst-export/v3.Exporter URI. (default "whosonfirst://")
  -float-property value
    	One or more {KEY}={VALUE} flags where {KEY} is a valid tidwall/gjson path and {VALUE} is a float(64) value.
  -int-property value
    	One or more {KEY}={VALUE} flag where {KEY} is a valid tidwall/gjson path and {VALUE} is a int(64) value.
  -iterator-uri string
    	A valid whosonfirst/go-whosonfirst-iterate/v3.Iterator URI. (default "repo://")
  -string-property value
    	One or more {KEY}={VALUE} flag where {KEY} is a valid tidwall/gjson path and {VALUE} is a string value.
  -verbose
    	Enable verbose (debug) logging
  -writer-uri string
    	A valid whosonfirst/go-writer/v3.Writer URI. (default "null://")
```

### Examples

For example:

```
$> ./bin/wof ensure-property \
	-boolean-property properties.wfdn:is_restricted=false \
	-iterator-uri 'repo://?include=properties.mz:is_current=1&properties.sfomuseum:post_security=-1' \
	-writer-uri repo:///usr/local/data/sfomuseum-data-wayfinding/ \
	/usr/local/data/sfomuseum-data-wayfinding/
	
2025/07/17 10:57:51 INFO Updated record path=191/466/628/9/1914666289.geojson
2025/07/17 10:57:51 INFO Updated record path=191/466/629/5/1914666295.geojson
2025/07/17 10:57:51 INFO Updated record path=191/466/629/9/1914666299.geojson
2025/07/17 10:57:51 INFO Updated record path=191/466/630/3/1914666303.geojson
2025/07/17 10:57:51 INFO Updated record path=191/466/630/5/1914666305.geojson
2025/07/17 10:57:51 INFO Updated record path=191/466/631/9/1914666319.geojson
2025/07/17 10:57:51 INFO Updated record path=191/466/637/1/1914666371.geojson
... and so on
```