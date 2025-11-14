# wof

## wof ensure-property

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

### Examples

For example:

```
$> ./bin/wof ensure-property \
	-boolean-property properties.wfdn:is_restricted=false \
	-iterator-uri 'repo://?include=properties.mz:is_current=1&include=properties.sfomuseum:post_security=-1' \
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

Or to assign the geometry property for a collection of records (parented by WOF ID `1729828959`):

```
$> ./bin/ensure-property \
	-iterator-uri 'repo://?include=properties.wof:parent_id=1729828959' \
	-geometry-property 'wkt=POINT(-122.37385309950533 37.61939390848619)' \
	-verbose \
	-writer-uri repo:///usr/local/data/sfomuseum-data-publicart \
	/usr/local/data/sfomuseum-data-publicart
```

Or to assign an integer property derived from the value of another existing property in the same document:

```
$> ./bin/wof ensure-property \
	-if-missing \
	-int-property properties.georef:lastmodified=0 \
	-int-property-from properties.wof:lastmodified \
	-iterator-uri 'repo://?include=properties.georef:depicted=.*' \
	-writer-uri repo:///usr/local/data/sfomuseum-data-collection \
	/usr/local/data/sfomuseum-data-collection
	
2025/11/14 11:11:27 INFO Updated record path=151/192/646/5/1511926465.geojson
2025/11/14 11:11:29 INFO Updated record path=176/268/724/3/1762687243.geojson
2025/11/14 11:11:30 INFO Updated record path=176/282/829/5/1762828295.geojson
...and so on
```