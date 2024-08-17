# wof

## wof show

```
$> ./bin/wof show -h
Command-line tool for serving GeoJSON features from an on-demand web server.
Usage:
	 ./bin/wof path(N) path(N)
Valid options are:
  -label value
    	Zero or more (GeoJSON Feature) properties to use to construct a label for a feature's popup menu when it is clicked on.
  -map-provider string
    	Valid options are: leaflet, protomaps (default "leaflet")
  -map-tile-uri string
    	A valid Leaflet tile layer URI. See documentation for special-case (interpolated tile) URIs. (default "https://tile.openstreetmap.org/{z}/{x}/{y}.png")
  -point-style string
    	A custom Leaflet style definition for point geometries. This may either be a JSON-encoded string or a path on disk.
  -port int
    	The port number to listen for requests on (on localhost). If 0 then a random port number will be chosen.
  -protomaps-theme string
    	A valid Protomaps theme label. (default "white")
  -style string
    	A custom Leaflet style definition for geometries. This may either be a JSON-encoded string or a path on disk.

If the only path as input is "-" then data will be read from STDIN.
```

### Examples

For example:

```
$> ./bin/wof show montreal.geojson 
Records are viewable at http://localhost:63675
```

This should automatically open a new window in your default browser like this:

![](../docs/images/wof-show-montreal.png)

As of this writing there is minimal styling and little to no interactivity. That may (or may not) change. Right now this tool is principally just to be able to look at one or more Who's On First features on a map.
