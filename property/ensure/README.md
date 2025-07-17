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
