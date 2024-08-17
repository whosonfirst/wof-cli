# wof

## emit

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

### Examples

#### Example (CSV)

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

#### Example (GeoParquet)

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

#### Example (SPR)

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

#### Example (FeatureCollection)

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
