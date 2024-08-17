# wof

## wof validate

```
$> ./bin/wof validate -h
Validate one or more Who's On First documents.
Usage:
	 ./bin/wof path(N) path(N)
```

### Examples

For example:

```
$> ./bin/wof validate test.geojson
2024/05/20 15:19:34 Failed to run 'validate' command, Failed to validate body for 'test.geojson', Failed to validate name, Failed to derive wof:name from body, Missing wof:name property
```

Or:

```
$> curl -s https://data.whosonfirst.org/102527513 | ./bin/wof validate -
2024/05/20 15:47:34 Failed to run 'validate' command, Failed to validate body for '-', Failed to validate EDTF, Failed to validate EDTF cessation date, Unrecognized EDTF string 'open' (Invalid or unsupported EDTF string)
```

If I download the record for SFO and manually change the `edtf:cessation` property from "open" to ".." and then run the tool again, everything validates.

```
$> ./bin/wof validate sfo.geojson
```
