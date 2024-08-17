# wof

## wof geometry

```
$> ./bin/wof geometry -h
Display or update the geometry for one or more Who's On First records.
Usage:
	 ./bin/wof geometry [options] path(N) path(N)
  -action string
    	Valid options are: show, update. (default "show")
  -format string
    	The format of the source geometry (if -action update). Valid options are: geojson. (default "geojson")
  -source string
    	The data source from which to derive a geometry (if -action update). Valid options are: the path to a local file.
  -stdout
    	Boolean flag signaling that updated records should be written to STDOUT. If false input files will be overwritten.
```

### Examples

For example:

```
$> wof geometry \
	-action update \
	-source /usr/local/data/t1.geojson \
	/usr/local/data/sfomuseum-data-architecture/data/191/456/434/5/1914564345.geojson
```

Or:

```
$> wof geometry 1914564345
{"coordinates":[[[[-122.385635,37.611836],[-122.385693,37.611861],[-122.385712,37.611868],[-122.385789,37.611901],[-122.385836,37.61192],[-122.385915,37.611954],[-122.385906,37.611969],[-122.385902,37.611974],[-122.385902,37.611974],[-122.385867,37.612026],[-122.385867,37.612026],[-122.385853,37.612048],[-122.385846,37.612046],[-122.385754,37.612007],[-122.385572,37.61193],[-122.385446,37.61212] ... and so on
```
