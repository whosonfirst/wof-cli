# wof-cli

Experimental 'wof' command-line tool for common Who's On First operations.

## Documentation

Documentation is incomplete at this time.

## Tools

```
> make cli
go build -mod vendor -ldflags="-s -w" -o bin/wof cmd/wof/main.go
```

### wof

```
$> ./bin/wof -h
Usage: wof [CMD] [OPTIONS]
Valid commands are:
* export
* fmt
* pip
* validate
```

_Important: The inputs and outputs for the `wof` tool have not been finalized yet, notably about how files are read and written if updated. You should expect change in the short-term._

#### wof export

```
$> ./bin/wof export -h
"Export-ify" one or more Who's On First records.
Usage:
	 ./bin/wof path(N) path(N)
  -stdout
    	Boolean flag signaling that updated records should be written to STDOUT. If false input files will be overwritten. (default true)
```

#### wof fmt

```
$> ./bin/wof fmt -h
Format one or more GeoJSON files according to Who's On First conventions.
Usage:
	 ./bin/wof path(N) path(N)
  -stdout
    	Boolean flag signaling that updated records should be written to STDOUT. If false input files will be overwritten. (default true)
```

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

#### wof pip

```
$> ./bin/wof pip -h
Perform point-in-polygon and wof:hierarchy update operations on one or more Who's On First records.
Usage:
	 ./bin/wof path(N) path(N)
  -export
    	"Export-ify" each record after point-in-polygon operations are complete.
  -mapshaper-client-uri string
    	Optional URI to a sfomuseum/go-sfomuseum-mapshaper server instance used to derive point-in-polygon centroids. If absent then the centroid used to perform point-in-polygon operations will be determined using internal heuristics.
  -placetype string
    	Assign this value as the "wof:placetype" property before performing point-in-polygon operations.
  -spatial-database-uri string
    	A valid whosonfirst/go-whosonfirst-spatial/database URI. By default 'pmtiles://' and 'sqlite://' spatial database URIs are supported.
  -stdout
    	Boolean flag signaling that updated records should be written to STDOUT. (default true)
```

For example, imagine a shell script like this (in order to account for boring URL-escaping issues):

```
#!/bin/sh

TILES_URI=file:///usr/local/data
TILES_DATABASE=architecture
TILES_ZOOM=12
TILES_LAYER=architecture

ENC_TILES_URI=`urlescape ${TILES_URI}`

./bin/wof pip \
    -spatial-database-uri "pmtiles://?tiles=${ENC_TILES_URI}&database=${TILES_DATABASE}&zoom=${TILES_ZOOM}&layer=${TILES_LAYER}&enable-cache=true" \
    -mapshaper-client-uri http://localhost:8080 \
    -placetype venue \
    test.geojson 
	 

```

Running this command would yield something like this:

```
2024/05/20 15:34:27 INFO fetching architecture 0-16384
2024/05/20 15:34:27 INFO fetched architecture 0-16384
2024/05/20 15:34:27 INFO fetching architecture 127-89
2024/05/20 15:34:27 INFO fetched architecture 127-89
2024/05/20 15:34:27 INFO Time to create database path=/architecture/12/655/1585.mvt "spatial database uri"="sqlite://?dsn=modernc://mem" "count features"=628 time=183.835708ms
{
  "id": 6102,
  "type": "Feature",
  "properties": {
    "USAGE": "STAIR"
  ,"wof:placetype":"venue","wof:parent_id":102527513,"wof:country":"US","wof:hierarchy":[{"campus_id":102527513,"continent_id":102191575,"country_id":85633793,"county_id":102087579,"locality_id":85922583,"postalcode_id":554784711,"region_id":85688637},{"campus_id":102527513,"continent_id":102191575,"country_id":85633793,"county_id":102085387,"region_id":85688637}]},
  "bbox": null,
  "geometry": {"coordinates":[[[-122.39515729496134,37.62304022215619,40],[-122.39511251529856,37.62305107289922,40],[-122.39508811697748,37.623056984354164,40],[-122.39509251970546,37.62306848116357,40],[-122.39510034575342,37.62308891418595,40],[-122.39516943471905,37.62307271147819,40],[-122.39515729496134,37.62304022215619,40]]],"type":"Polygon"}
}
```

There are a few things to note:

* The `-placetype` flag is a convenience to facilitate point-in-polygon operations without having to first update an input record.
* By default the `pip` command neither "formats" or "exports" results. Although there is an `-export` flag to enable both it is your responsibility to ensure that input documents have all the necessary properties (for example "wof:name").
* The `pip` command does _NOT_ yet implement the logic of the [py-mapzen-whosonfirst-hierarchy](https://github.com/whosonfirst/py-mapzen-whosonfirst-hierarchy/blob/master/mapzen/whosonfirst/hierarchy/__init__.py) library. There is an [open issue](https://github.com/whosonfirst/go-whosonfirst-spatial/issues/40) for this.

#### wof validate

```
$> ./bin/wof validate -h
Validate one or more Who's On First documents.
Usage:
	 ./bin/wof path(N) path(N)
```

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

_Note the default behaviour for successfull validation is to do nothing. That might change? Maybe?_

## See also

* https://github.com/whosonfirst/go-whosonfirst-export
* https://github.com/whosonfirst/go-whosonfirst-format
* https://github.com/whosonfirst/go-whosonfirst-validate
* https://github.com/whosonfirst/go-whosonfirst-spatial