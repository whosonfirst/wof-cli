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
* emit
* export
* fmt
* open
* pip
* show
* validate
```

_Important: The inputs and outputs for the `wof` tool have not been finalized yet, notably about how files are read and written if updated. You should expect change in the short-term._

#### wof emit

```
$> ./bin/wof emit -h
Emit one or more Who's On First records.
Usage:
	 ./bin/wof emit [options] path(N) path(N)
  -as-spr
    	Emit Who's On First records formatted as Standard Place Response (SPR) records.
  -as-spr-geojson
    	Emit Who's On First records as GeoJSON records where the 'properties' element is replaced by a Standard Place Response (SPR) representation of the record.
  -iterator-uri string
    	A valid whosonfirst/go-whosonfirst-iterate/v2/emitter URI. (default "repo://")
  -writer-uri string
    	A valid whosonfirst/go-writer.Writer URI. (default "jsonl://?writer=stdout://")
```

For example:

```
$> ./bin/wof emit -as-spr -writer-uri 'jsonl://?writer=stdout://' /usr/local/data/sfomuseum-data-maps/ 
{"edtf:cessation":"1985~","edtf:inception":"1985~","mz:is_ceased":1,"mz:is_current":0,"mz:is_deprecated":-1,"mz:is_superseded":0,"mz:is_superseding":0,"mz:latitude":37.616459,"mz:longitude":-122.386272,"mz:max_latitude":37.63100646804649,"mz:max_longitude":-122.37094362769881,"mz:min_latitude":37.60096637420677,"mz:min_longitude":-122.40407820844655,"mz:uri":"https://data.whosonfirst.org/136/039/131/3/1360391313.geojson","wof:belongsto":[],"wof:country":"US","wof:id":1360391313,"wof:lastmodified":1716594274,"wof:name":"SFO (1985)","wof:parent_id":-4,"wof:path":"136/039/131/3/1360391313.geojson","wof:placetype":"custom","wof:repo":"sfomuseum-data-maps","wof:superseded_by":[],"wof:supersedes":[]}
...and so on
```

Or:

```
$> ./bin/wof emit -as-spr -as-spr-geojson \
	-writer-uri 'featurecollection://?writer=stdout://' \
	/usr/local/data/sfomuseum-data-flights-2024-05 \
	
| ogr2ogr -f parquet flights2.parquet /vsistdin/
```

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

#### wof open

```
$> ./bin/wof open -h
Open one or more Who's On First documents in a custom editor.
Usage:
	 ./bin/wof path(N) path(N)
  -editor string
    	The editor to use for opening Who's On First records. If empty the value of the EDITOR environment variable will be used.
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

#### wof show

```
$> ./bin/wof show -h
Display one or more Who's On First documents on a map.
Usage:
	 ./bin/wof path(N) path(N)
  -port int
    	The port number to listen for requests on (on localhost). If 0 then a random port number will be chosen.
```

For example:

```
$> ./bin/wof show montreal.geojson 
Records are viewable at http://localhost:63675
```

This should automatically open a new window in your default browser like this:

![](docs/images/wof-show-montreal.png)

As of this writing there is minimal styling and little to no interactivity. That may (or may not) change. Right now this tool is principally just to be able to look at one or more Who's On First features on a map.

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

## Paths and URIs

The default behaviour for the `wof` command is to assume that all the URIs it is passed are paths on the local computer.

That being said all the tool also "expand" each path using the [uris.ExpandURIsWithCallback](uris/uris.go) method which allows for more sophisticated behaviour in the future.

Currently there is only one supported "expansion": Treating bare numbers as Who's On First IDs, resolving them to their relative path and looking for that file in a parent "data" directory inside the current working directory. Basically it's a shortcut for resolving a Who's On First record to its GeoJSON representation inside a `whosonfirst-data` repository.

Eventually there may be other "expansions" most notably support for the Go `./...` syntax to process all the Who's On First records in the current working directory.

## See also

* https://github.com/whosonfirst/go-whosonfirst-export
* https://github.com/whosonfirst/go-whosonfirst-format
* https://github.com/whosonfirst/go-whosonfirst-format-wasm
* https://github.com/whosonfirst/go-whosonfirst-validate
* https://github.com/whosonfirst/go-whosonfirst-spatial
* https://github.com/whosonfirst/go-whosonfirst-spatial-pmtiles
* https://github.com/whosonfirst/go-whosonfirst-spatial-sqlite
* https://github.com/sfomuseum/go-sfomuseum-mapshaper