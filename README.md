# wof-cli

Command-line tool for common Who's On First operations.

## Tools

```
$> make cli
go build -mod vendor -ldflags="-s -w" -o bin/wof cmd/wof/main.go

Features and functionality are enable through the use of (Go language) [build tags](#). By default all build tags are enabled. Please consult the [Build tags](#) section below for detailed descriptions.

### wof

```
$> ./bin/wof -h
Usage: wof [CMD] [OPTIONS]
Valid commands are:
* centroid
* emit
* export
* fmt
* geometry
* open
* pip
* property
* show
* supersede
* uri
* validate
```

_Important: The inputs and outputs for the `wof` tool have not been finalized yet, notably about how files are read and written if updated. You should expect change in the short-term._

#### wof centroid

```
$> ./bin/wof centroid -h
Emit centroids data (source, latitude, longitude) for one or more Who's On First records as CSV-encoded properties.
Usage:
	 ./bin/wof emit [options] path(N) path(N)
```

For example:

```
$> /usr/local/whosonfirst/wof-cli/bin/wof centroid 102536223
uri,source,latitude,longitude
/usr/local/data/sfomuseum-data-whosonfirst/data/102/536/223/102536223.geojson,lbl,21.040317,-86.872856
```

#### wof emit

```
$> ./bin/wof emit -h
Emit one or more Who's On First records.
Usage:
	 ./bin/wof emit [options] path(N) path(N)
  -as-spr
    	Emit Who's On First records formatted as Standard Place Response (SPR) records. This flag is DEPRECATED. Please use '-format spr' instead.
  -as-spr-geojson
    	Emit Who's On First records as GeoJSON records where the 'properties' element is replaced by a Standard Place Response (SPR) representation of the record. This flag is DEPRECATED. Please use '-format geojson' instead.
  -csv-append-property value
    	Zero or more additional properties to append to each CSV row. Properties should be in the format of {COLUMN_NAME}={PATH}. This flag is only honoured if the -format flag has a value of "csv".
  -forgiving
    	Do not stop processing when errors are encountered.
  -format string
    	Valid options are: csv, spr, spr-geojson or [none]. If none then the raw GeoJSON for each matching record will be emitted.
  -include-alt-geoms
    	Emit alternate geometry records. (default true)
  -iterator-uri string
    	A valid whosonfirst/go-whosonfirst-iterate/v2/emitter URI. If URI is "-" then this flag will be assigned a value of "file://" whose input will be the expanded URIs derived from additional arguments. Available options are: cwd://, directory://, featurecollection://, file://, filelist://, geojsonl://, null://, repo:// (default "repo://")
  -query value
    	One or more {PATH}={REGEXP} parameters for filtering records.
  -query-mode string
    	Specify how query filtering should be evaluated. Valid modes are: ALL, ANY (default "ALL")
  -writer-uri string
    	A valid whosonfirst/go-writer.Writer URI. Available options are: cwd://, featurecollection://, fs://, geoparquet://, io://, jsonl://, null://, repo://, sqlite://, stdout:// (default "jsonl://?writer=stdout://")
```

##### Example (CSV)

For example, emitting records as CSV results with additional custom properties:

```
$> ./bin/wof emit \
	-format csv \
	-csv-append-property sfomuseum_description=properties.sfomuseum:description \
	-writer-uri stdout:// \
	-query 'properties.millsfield:subcategory_id=1511213363' \
	/usr/local/data/sfomuseum-data-collection

wof:repo,mz:is_deprecated,edtf:inception,mz:max_latitude,mz:is_ceased,mz:is_superseding,wof:lastmodified,mz:longitude,mz:max_longitude,wof:placetype,wof:country,mz:is_current,wof:belongsto,mz:latitude,mz:min_longitude,wof:superseded_by,mz:min_latitude,edtf:cessation,wof:parent_id,wof:path,mz:is_superseded,wof:name,wof:id,mz:uri,wof:supersedes,sfomuseum_description
sfomuseum-data-collection,-1,1948,37.61661101879963,1,0,1716579582,-122.38615540108617,-122.38583385944366,custom,XY,1,"102527513,1159162825,102191575,85633793,102087579,85922583,85688637,1159160869",37.61635747477365,-122.38647222518921,,37.6161053594541,1948,1159162825,184/675/460/1/1846754601.geojson,0,postcard: Pigeon Key,1846754601,184/675/460/1/1846754601.geojson,,"Color postcard with photographic image depicting houses on island with bridge; postmarked June 7, 1948 in Fort Lauderdale, Fla.; text on front: “Kodachrome by C.H. Ruth / Overseas Highway, Above Pigeon Key, Between Key West and Miami, Fla.”"
Black and white postcard with photographic image depicting profile view of Sikorsky XPBS-1 on water; text on front: “Sikorsky flying dreadnaught [sic] for U/S/ Navy with Hamilton Constant Speed Propellers”.,custom,184/675/460/3/1846754603.geojson,XY,37.6161053594541,1,1716579582,1159162825,0,1930-12,184/675/460/3/1846754603.geojson,"102527513,1159162825,102191575,85633793,102087579,85922583,85688637,1159160869",-1,,37.61661101879963,1930-09,,1,-122.38615540108617,-122.38647222518921,sfomuseum-data-collection,0,postcard: Sikorsky XPBS-1,-122.38583385944366,1846754603,37.61635747477365
XY,1930~,"Black and white postcard with photographic image depicting aerial view of Havana Harbor with Pan American Airways Fokker F-10 in flight; text on front: “Tri-motor airliner, Pan American Airways, over El Morro, Havana”.",1716579582,-122.38583385944366,1,0,184/675/460/5/1846754605.geojson,0,37.61661101879963,-122.38647222518921,37.61635747477365,"102527513,1159162825,102191575,85633793,102087579,85922583,85688637,1159160869",-122.38615540108617,"postcard: Pan American Airways, Fokker F.10, Havana",custom,1159162825,sfomuseum-data-collection,1930~,,37.6161053594541,,1846754605,184/675/460/5/1846754605.geojson,1,-1
-122.38583385944366,,37.61661101879963,193X,0,"102527513,1159162825,102191575,85633793,102087579,85922583,85688637,1159160869",0,184/675/460/7/1846754607.geojson,1159162825,193X,sfomuseum-data-collection,"Black and white postcard with photographic image depicting low-angle front three-quarter view of Pan American Airways Sikorsky S-42 in flight against clouds; text on front: “A ‘Clipper Ship’ of the Air”; text on reverse: “America’s Largest Airliners - Giant 4-engined 32 and 44 passenger, 19-ton Flying Boats - Ply the Pan American Airways Routes between the U.S., the West Indies and South America”.","postcard: Pan American Airways, Sikorsky S-42",custom,37.61635747477365,,1,184/675/460/7/1846754607.geojson,-122.38647222518921,-1,1716579582,1846754607,-122.38615540108617,1,XY,37.6161053594541
... and so on
```

The default set of CSV row map to the properties of a Standard Places Result (SPR).

##### Example (GeoParquet)

For example, emitting all the records marked `mz:is_current=1` from the [whosonfirst-data-venue-ca](https://github.com/whosonfirst-data/whosonfirst-data-venue-ca) repository to a GeoParquet database:

```
./bin/wof emit \
	-writer-uri 'geoparquet://?min=100&max=1000&append-property=wof:concordances' \
	-iterator-uri 'repo://?include=properties.mz:is_current=1' \
	/usr/local/data/whosonfirst-data-venue-ca \
	> /usr/local/data/venue-ca.geoparquet

processed 17349 records in 1m0.000336833s (started 2024-08-16 17:38:11.279896 -0700 PDT m=+0.044642876)
... time passes
processed 546529 records in 20m0.002204708s (started 2024-08-16 17:38:11.279896 -0700 PDT m=+0.044642876)
2024/08/16 17:59:07 INFO time to index paths (1) 20m56.026484291s
```

And then loading the resultant database in [DuckDB](https://duckdb.org/):

```
$> duckdb
v1.0.0 1f98600c2c
Enter ".help" for usage hints.
Connected to a transient in-memory database.
Use ".open FILENAME" to reopen on a persistent database.
D LOAD spatial;
D SELECT "wof:id", "wof:name", "wof:concordances" FROM read_parquet('/usr/local/data/venue-ca.geoparquet') LIMIT 5;
┌────────────┬───────────────────┬──────────────────────────────────────────────────────────────────────────────────────────────┐
│   wof:id   │     wof:name      │                                       wof:concordances                                       │
│  varchar   │      varchar      │                                           varchar                                            │
├────────────┼───────────────────┼──────────────────────────────────────────────────────────────────────────────────────────────┤
│ 1108798699 │ Foxy              │ {\n      "osm:node": 2687914646,\n      "ovtr:id": "08f2baa44c8f6028031206a8c92b602c"\n    } │
│ 1108798581 │ Myriade           │                                                                                              │
│ 1125142651 │ Liverpool House   │ {}                                                                                           │
│ 1125142781 │ Joe Beef          │ {\n      "4sq:id": "4b84801af964a520fd3831e3"\n    }                                         │
│ 1108808935 │ Drawn & Quarterly │ {\n      "4sq:id": "4ad4c06ff964a5205ffb20e3",\n      "osm:node": 2704382357\n    }          │
└────────────┴───────────────────┴──────────────────────────────────────────────────────────────────────────────────────────────┘
D
```

Or reading data directly from a [Who's On First -style data repository on GitHub](https://github.com/sfomuseum-data/sfomuseum-data-architecture) and writing all the records to GeoParquet file:

```
$> ./bin/wof emit \
	-writer-uri 'geoparquet://?min=100&max=1000&append-property=sfomuseum:placetype' \
	-iterator-uri 'git:///tmp' \
	https://github.com/sfomuseum-data/sfomuseum-data-architecture.git \
	> /usr/local/data/arch.geoparquet
```

And then, again, in DuckDB:

```
D SELECT "wof:id", "wof:name", "sfomuseum:placetype", "wof:placetype", "mz:is_current" FROM read_parquet('/usr/local/data/arch.geoparquet');
┌────────────┬─────────────────────────────────────────┬─────────────────────┬───────────────┬───────────────┐
│   wof:id   │                wof:name                 │ sfomuseum:placetype │ wof:placetype │ mz:is_current │
│  varchar   │                 varchar                 │       varchar       │    varchar    │    varchar    │
├────────────┼─────────────────────────────────────────┼─────────────────────┼───────────────┼───────────────┤
│ 102527513  │ San Francisco International Airport     │ airport             │ campus        │ 1             │
│ 1159157037 │ G-02 International North Cases          │ gallery             │ enclosure     │ 0             │
│ 1159157039 │ A-02 International South Cases          │ gallery             │ enclosure     │ 0             │
│ 1159157041 │ A-07 International Central Vitrine      │ gallery             │ enclosure     │ 0             │
│ 1159157045 │ 4B International North Wall             │ gallery             │ enclosure     │ 0             │
│ 1159157047 │ 4C International South Wall             │ gallery             │ enclosure     │ 0             │
│ 1159157049 │ 3L Terminal 3 Connector Arrival Level   │ gallery             │ enclosure     │ 0             │
│ 1159157051 │ 3J Photographs                          │ gallery             │ enclosure     │ 0             │
│ 1159157053 │ 3C North Connector                      │ gallery             │ enclosure     │ 0             │
│ 1159157055 │ 3D Terminal 3 Hub                       │ gallery             │ enclosure     │ 0             │
... and so on
```

Please consult the [whosonfirst/go-writer-geoparquet](https://github.com/whosonfirst/go-writer-geoparquet?tab=readme-ov-file#how-does-it-work) documentation for details on how to configure the `-writer-uri` flag.

##### Example (SPR)

For example, emitting records as SPR results:

```
$> ./bin/wof emit -as-spr -writer-uri 'jsonl://?writer=stdout://' /usr/local/data/sfomuseum-data-maps/ 

{"edtf:cessation":"1985~","edtf:inception":"1985~","mz:is_ceased":1,"mz:is_current":0,"mz:is_deprecated":-1,"mz:is_superseded":0,"mz:is_superseding":0,"mz:latitude":37.616459,"mz:longitude":-122.386272,"mz:max_latitude":37.63100646804649,"mz:max_longitude":-122.37094362769881,"mz:min_latitude":37.60096637420677,"mz:min_longitude":-122.40407820844655,"mz:uri":"https://data.whosonfirst.org/136/039/131/3/1360391313.geojson","wof:belongsto":[],"wof:country":"US","wof:id":1360391313,"wof:lastmodified":1716594274,"wof:name":"SFO (1985)","wof:parent_id":-4,"wof:path":"136/039/131/3/1360391313.geojson","wof:placetype":"custom","wof:repo":"sfomuseum-data-maps","wof:superseded_by":[],"wof:supersedes":[]}
...and so on
```

Or with query filtering:

```
$> ./bin/wof emit -as-spr -query 'properties.wof:name=SFO \(2023\)' /usr/local/data/sfomuseum-data-maps/

{"edtf:cessation":"","edtf:inception":"2023-07~","mz:is_ceased":-1,"mz:is_current":-1,"mz:is_deprecated":-1,"mz:is_superseded":0,"mz:is_superseding":0,"mz:latitude":37.621284127293116,"mz:longitude":-122.38285759138246,"mz:max_latitude":37.642285759714994,"mz:max_longitude":-122.34578162574567,"mz:min_latitude":37.60153229886917,"mz:min_longitude":-122.40810153962025,"mz:uri":"https://data.whosonfirst.org/188/030/951/9/1880309519.geojson","wof:belongsto":[102527513,102191575,85633793,102087579,85922583,554784711,85688637,102085387],"wof:country":"US","wof:id":1880309519,"wof:lastmodified":1716594274,"wof:name":"SFO (2023)","wof:parent_id":-1,"wof:path":"188/030/951/9/1880309519.geojson","wof:placetype":"custom","wof:repo":"sfomuseum-data-maps","wof:superseded_by":[],"wof:supersedes":[]}
```

##### Example (FeatureCollection)

Or emitting records as FeatureCollection of GeoJSON-formatted SPR results (where the original geometry is preserved but the properties hash is replaced by that record's SPR) and piping the result to `ogr2ogr`:

```
$> ./bin/wof emit -as-spr -as-spr-geojson \
	-writer-uri 'featurecollection://?writer=stdout://' \
	/usr/local/data/sfomuseum-data-flights-2024-05 \
	
| ogr2ogr -f parquet flights2.parquet /vsistdin/
```

Or, iterating over a custom list of files:

```
$> wof emit -iterator-uri - -writer-uri 'featurecollection://?writer=stdout://' 1914563993 1914564157 1914564489 1914564345 | json_pp | grep 'wof:name'
2024/07/09 13:57:07 INFO time to index paths (4) 797.084µs
            "wof:name" : "Terminal 2",
            "wof:name" : "Terminal 1",
            "wof:name" : "International Terminal",
            "wof:name" : "Terminal 3",
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

_Note: Exporting alternate geometry files is not supported yet._

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

#### wof geometry

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

For example:

```
$> wof geometry -action update -source /usr/local/data/t1.geojson /usr/local/data/sfomuseum-data-architecture/data/191/456/434/5/1914564345.geojson
```

Or:

```
$> wof geometry 1914564345
{"coordinates":[[[[-122.385635,37.611836],[-122.385693,37.611861],[-122.385712,37.611868],[-122.385789,37.611901],[-122.385836,37.61192],[-122.385915,37.611954],[-122.385906,37.611969],[-122.385902,37.611974],[-122.385902,37.611974],[-122.385867,37.612026],[-122.385867,37.612026],[-122.385853,37.612048],[-122.385846,37.612046],[-122.385754,37.612007],[-122.385572,37.61193],[-122.385446,37.61212] ... and so on
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

#### wof property

Print one or more properties for one or more Who's On First IDs.

```
$> wof property -h
Print one or more properties for one or more Who's On First IDs.
Usage:
	 wof path(N) path(N)
  -format string
    	Valid options are: csv. If empty then properties will printed as a new-line separated list.
  -path value
    	One or more valid tidwall/gjson paths to extract from each document
  -prefix string
    	If not empty this prefix will be appended (and separated by a ".") to each -path argument
```

For example:

```
$> wof property -path properties.wof:name 1796903597 1796889561 1796889543 1796889557 1796903629 1796889563 1796935715 1796935615
AirTrain Gargage G / BART Red Line
AirTrain Long-Term Parking Blue Line (Outbound)
AirTrain Garage G / BART Blue Line
AirTrain Westfield Road Station (Inbound)
Grand Hyatt Hotel Reception
Rental Car Center
International Terminal Main Hall Departures Door 2
San Francisco International Airport BART Station Platform
```

It is also possible to emit properties for records as CSV data by passing the `-format csv` flag. For example:

```
$> wof property -format csv -path properties.wof:name -path properties.wof:parent_id 1796903597 1796889561 1796889543 1796889557 1796903629 1796889563 1796935715 1796935615
uri,properties.wof:name,properties.wof:parent_id
/usr/local/data/sfomuseum-data-wayfinding/data/179/690/359/7/1796903597.geojson,AirTrain Gargage G / BART Red Line,1477855991
/usr/local/data/sfomuseum-data-wayfinding/data/179/688/956/1/1796889561.geojson,AirTrain Long-Term Parking Blue Line (Outbound),102527513
/usr/local/data/sfomuseum-data-wayfinding/data/179/688/954/3/1796889543.geojson,AirTrain Garage G / BART Blue Line,1477855991
/usr/local/data/sfomuseum-data-wayfinding/data/179/688/955/7/1796889557.geojson,AirTrain Westfield Road Station (Inbound),102527513
/usr/local/data/sfomuseum-data-wayfinding/data/179/690/362/9/1796903629.geojson,Grand Hyatt Hotel Reception,1477856005
/usr/local/data/sfomuseum-data-wayfinding/data/179/688/956/3/1796889563.geojson,Rental Car Center,1477863277
/usr/local/data/sfomuseum-data-wayfinding/data/179/693/571/5/1796935715.geojson,International Terminal Main Hall Departures Door 2,1745882445
/usr/local/data/sfomuseum-data-wayfinding/data/179/693/561/5/1796935615.geojson,San Francisco International Airport BART Station Platform,102527513
```

CSV-formatted output will automatically append a `uri` column to each row.

If you know that all the `-path` flags share the same prefix you can specify it using the `-prefix` flag and save the time it will take you to include it with each `-path` flag. For example:

```
$> wof property -format csv -prefix properties -path wof:name -path wof:superseded_by 1477855991 102527513 1477855991 102527513 1477856005 1477863277 1745882445 102527513
uri,properties.wof:name,properties.wof:superseded_by
/usr/local/data/sfomuseum-data-architecture/data/147/785/599/1/1477855991.geojson,Garage G,[]
/usr/local/data/sfomuseum-data-architecture/data/102/527/513/102527513.geojson,San Francisco International Airport,[]
/usr/local/data/sfomuseum-data-architecture/data/147/785/599/1/1477855991.geojson,Garage G,[]
/usr/local/data/sfomuseum-data-architecture/data/102/527/513/102527513.geojson,San Francisco International Airport,[]
/usr/local/data/sfomuseum-data-architecture/data/147/785/600/5/1477856005.geojson,Grand Hyatt Hotel,[]
/usr/local/data/sfomuseum-data-architecture/data/147/786/327/7/1477863277.geojson,Rental Car Center,[]
/usr/local/data/sfomuseum-data-architecture/data/174/588/244/5/1745882445.geojson,International Terminal Arrivals,[]
/usr/local/data/sfomuseum-data-architecture/data/102/527/513/102527513.geojson,San Francisco International Airport,[]
```

#### wof show

```
> ./bin/wof show -h
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

For example:

```
$> ./bin/wof show montreal.geojson 
Records are viewable at http://localhost:63675
```

This should automatically open a new window in your default browser like this:

![](docs/images/wof-show-montreal.png)

As of this writing there is minimal styling and little to no interactivity. That may (or may not) change. Right now this tool is principally just to be able to look at one or more Who's On First features on a map.

#### wof supersede

Supersede one or more Who's On First record and update the records superseding them.

```
$> ./bin/wof supersede -h
Supersede one or more Who's On First record and update the records superseding them.
Usage:
	 ./bin/wof [options] path(N) path(N)
  -deprecated
    	Each record being superseded should also be marked as deprecated.
  -parent-id int
    	The parent ID to assign to the new record. (default -1)
  -parent-reader-uri -1
    	A valid whosonfirst/go-reader URI used to load parent records if -superseded-id is -1. Required if -parent-id is not `-1`. (default "null://")
  -superseding-id -1
    	The ID to supersede each record with. If -1 then each record will be cloned and the new ID of the clone will be used as the superseding_id. (default -1)
  -superseding-reader-uri -1
    	A valid whosonfirst/go-reader URI used to load records that are doing the superseding. Required if -superseding-id is not -1. (default "null://")
  -superseding-writer-uri -1
    	A valid whosonfirst/go-writer URI used to update records that are doing the superseding. Required if -superseding-id is not -1. (default "null://")
```

For example, supersede a set of records making the new records "clones" of the old records:

```
$> wof supersede \
	-superseding-writer-uri repo:///usr/local/data/sfomuseum-data-wayfinding \
	1796889561 1796889543 1796889557 1796903629 1796889563 1796935715 1796935615
```

_Important: As of this writing this tool does NOT supersede alternate geometries associated with the WOF records being superseded._

##### "reader" and "writer" URIs

The `wof` tool has its own internal logic for [deriving paths for reading and writing input documents](https://github.com/whosonfirst/wof-cli?tab=readme-ov-file#paths-and-uris).

That being the case by the time a document URI is resolved at the command layer there is not necessarily enough information to write documents _related_ to the document currently being processed. Further even if that context is known it may not be appropriate. For example, if a document in the `sfomuseum-data-wayfinding` repository is being superseded and the new document is parented by a document in the `sfomuseum-data-architecture` repository (using the `-parent-id` flag) then there is nothing in the `sfomuseum-data-wayfinding` context to know where to find that record.

Hence the `-parent-reader-uri`, `-superseding-reader-uri` and `-superseding-writer-uri` flags. There are expected to be valid [whosonfirst/go-reader.Reader](https://github.com/whosonfirst/go-reader) and [whosonfirst/go-writer.Writer](https://github.com/whosonfirst/go-writer) URIs. For example:

```
$> wof supersede \
	-parent-id 102527513 \
	-parent-reader-uri repo:///usr/local/data/sfomuseum-data-architecture \
	-superseding-writer-uri repo:///usr/local/data/sfomuseum-data-wayfinding \	
	1796889561 
```

It's a bit cumbersome but the decision was taken, given the potential for many unrelated moving parts, to be explicit rather than clever.

#### wof uri

```
$> ./bin/wof uri -h
Print the nested URI for one or more Who's On First IDs.
Usage:
	 ./bin/wof path(N) path(N)
  -prefix string
    	An optional prefix to append to the final URI.
```

For example:

```
$> ./bin/wof uri 1914650585
191/465/058/5/1914650585.geojson

$> cat `wof uri -prefix data 1914650737` | jq '.properties["wof:name"]'
"1H Student Art North"
```

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

### Expansions

Currently there is only two supported "expansions":

1. Treating bare numbers as Who's On First IDs, resolving them to their relative path and looking for that file in a parent "data" directory inside the current working directory. Basically it's a shortcut for resolving a Who's On First record to its GeoJSON representation inside a `whosonfirst-data` repository.
2. URIs prefixed with `repo://` will be treated as a [whosonfirst/go-whosonfirst-iterate/v2 "repo" emitter URI](https://github.com/whosonfirst/go-whosonfirst-iterate?tab=readme-ov-file#repo) and all the files in that repository will be processed.

Eventually there may be other "expansions" most notably support for the Go `./...` syntax to process all the Who's On First records in the current working directory.

## Advanced

### Build tags

For example to build the `wof` command line tool without suport for the [whosonfirst/go-writer-geoparquet](https://github.com/whosonfirst/go-writer-geoparquet) package (which would shave about 20MB off the size of the final binary) you would do this:

```
$> go build -mod vendor -ldflags="-s -w" -tags no_writer_geoparquet -o bin/wof cmd/wof/main.go
```

#### Available build tags

##### no_centroid

Disable the `wof centroid` command.

##### no_emit

Disable the `wof emit` command.

##### no_export

Disable the `wof export` command.

##### no_format

Disable the `wof format` command.

##### no_geometry

Disable the `wof geometry` command.

##### no_iterator_git

Disable import of the [whosonfirst/go-whosonfirst-iterate-git/v2](https://github.com/whosonfirst/go-whosonfirst-iterate-git) package.

This affects the `wof emit` command.

##### no_iterator_org

Disable import of the [whosonfirst/go-whosonfirst-iterate-organization/v2](https://github.com/whosonfirst/go-whosonfirst-iterate-organization) package.

This affects the `wof emit` command.

##### no_open

Disable the `wof open` command.

##### no_pip

Disable the `wof pip` command.

##### no_property

Disable the `wof property` command.

##### no_show

Disable the `wof show` command.

##### no_spatial_pmtiles

Disable import of the [whosonfirst/go-whosonfirst-spatial-pmtiles](https://github.com/whosonfirst/go-whosonfirst-spatial-pmtiles) package.

This affects the `wof pip` command.

##### no_spatial_sqlite

Disable import of the [whosonfirst/go-whosonfirst-spatial-sqlite](https://github.com/whosonfirst/go-whosonfirst-spatial-sqlite) package.

This affects the `wof pip` command.

##### no_supersede

Disable the `wof supersede` command.

##### no_uri

Disable the `wof uri` command.

##### no_validate

Disable the `wof validate` command.

##### no_writer_geoparquet

Disable import of the [whosonfirst/go-writer-geoparquet](https://github.com/whosonfirst/go-writer-geoparquet) package.

This affects the `wof emit` command.

## See also

* https://github.com/whosonfirst/go-whosonfirst-export
* https://github.com/whosonfirst/go-whosonfirst-format
* https://github.com/whosonfirst/go-whosonfirst-format-wasm
* https://github.com/whosonfirst/go-whosonfirst-validate
* https://github.com/whosonfirst/go-whosonfirst-spatial
* https://github.com/whosonfirst/go-whosonfirst-spatial-pmtiles
* https://github.com/whosonfirst/go-whosonfirst-spatial-sqlite
* https://github.com/whosonfirst/go-whosonfirst-iterate
* https://github.com/sfomuseum/go-sfomuseum-mapshaper
