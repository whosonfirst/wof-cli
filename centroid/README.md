# wof

## centroid

```
$> ./bin/wof centroid -h
Emit centroids data (source, latitude, longitude) for one or more Who's On First records as CSV-encoded properties.
Usage:
	 ./bin/wof emit [options] path(N) path(N)
```

### Examples

For example:

```
$> /usr/local/whosonfirst/wof-cli/bin/wof centroid 102536223
uri,source,latitude,longitude
/usr/local/data/sfomuseum-data-whosonfirst/data/102/536/223/102536223.geojson,lbl,21.040317,-86.872856
```