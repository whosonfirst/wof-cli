# wof

## wof uri

```
$> ./bin/wof uri -h
Print the nested URI for one or more Who's On First IDs.
Usage:
	 ./bin/wof path(N) path(N)
  -prefix string
    	An optional prefix to append to the final URI.
```

### Examples

For example:

```
$> ./bin/wof uri 1914650585
191/465/058/5/1914650585.geojson

$> cat `wof uri -prefix data 1914650737` | jq '.properties["wof:name"]'
"1H Student Art North"
```
