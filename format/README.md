# wof

## wof fmt

```
$> ./bin/wof fmt -h
Format one or more GeoJSON files according to Who's On First conventions.
Usage:
	 ./bin/wof path(N) path(N)
  -stdout
    	Boolean flag signaling that updated records should be written to STDOUT. If false input files will be overwritten. (default true)
```

### Examples

For example:

```
$> ./bin/wof fmt ~/Desktop/test.geojson
{
  "id": 6102,
  "type": "Feature",
  "properties": {
    "ACCESS_": "PUBLIC",
    "USAGE": "STAIR"
  },
  "bbox": null,
  "geometry": {"coordinates":[[[-122.39515729496134,37.62304022215619,40],[-122.39511251529856,37.62305107289922,40],[-122.39508811697748,37.623056984354164,40],[-122.39509251970546,37.62306848116357,40],[-122.39510034575342,37.62308891418595,40],[-122.39516943471905,37.62307271147819,40],[-122.39515729496134,37.62304022215619,40]]],"type":"Polygon"}
}
```
