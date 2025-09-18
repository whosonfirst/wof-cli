# wof cli cookbook

## Fetching all the `whosonfirst-data-admin-*` repositories as GeoParquet files (ogr2ogr)

```
#!/bin/sh

REPOS=`./bin/wof repos -prefix 'whosonfirst-data-admin-'`

for REPO in ${REPOS}
do

    echo "Fetch ${REPO}"

    if [ -f "${REPO}.parquet" ]
    then
	echo "Export exists for ${REPO}"
	continue
    fi
    
    echo "./bin/wof emit -format geojson -iterator-uri git:///tmp https://github.com/whosonfirst-data/${REPO}.git > ${REPO}.geojsonl"
    ./bin/wof emit -format geojson -iterator-uri git:///tmp https://github.com/whosonfirst-data/${REPO}.git > ${REPO}.geojsonl

    echo "ogr2ogr -f Parquet ${REPO}.parquet -oo OGRGeoJSONAllowNonStandard=YES 'GeoJSONSeq:${REPO}.geojsonl'"
    ogr2ogr -f Parquet ${REPO}.parquet -oo OGRGeoJSONAllowNonStandard=YES "GeoJSONSeq:${REPO}.geojsonl"

    echo "Clean up ${REPO}"
    rm ${REPO}.geojsonl
    
done
```

### Notes

* You could skip the 'wof repos' step entirely and simply do this `./bin/wof emit -format geojson -iterator-uri githuborg:///tmp 'whosonfirst-data://?prefix=whosonfirst-data-admin-'` but that will take a very long time and will not be terribly forgiving of errors (network, etc.)

* If you have GDAL installed with VSI support you can skip writing the `geojsonl` file simply do this: `./bin/wof emit -format geojson -iterator-uri git:///tmp https://github.com/whosonfirst-data/${REPO}.git | ogr2ogr -f Parquery ${REPO}.parquet -oo OGRGeoJSONAllowNonStandard=YES 'GeoJSONSeq:/vsisstdin/'`

* Because GeoParquet files are being created on a per-repo basis there is still the chance that there will be schema mismatches. There _shouldn't_ be but the reality is that it still happens, unfortunately. You could modify the script above to stream (append) all `.geojsonl` data to a single file and the perform the `ogr2ogr` transformation as a final step (outside of the `for REPO in ${REPOS}` loop).

## Fetching all the `whosonfirst-data-admin-*` repositories as GeoParquet files (native)

It is also possible to create a GeoParquet file from one or more Who's On First repos natively (as in "not needing to install ogr2ogr (gdal)"). For example:

```
TO_EMIT=""

REPOS=`./bin/wof repos -prefix 'whosonfirst-data-admin-'`

for $REPO in ${REPOS}
do
	TO_EMIT="${TO_EMIT} https://github.com/whosonfirst-data/${REPO}.git"
done

./bin/wof emit -writer-uri 'geoparquet://?writer=stdout://' -iterator-uri git:///tmp ${TO_EMIT} > wof.parquet
```

#### Notes

* Remember, this example is building a GeoParquet file for _all_ the Who's On First repositories so it will take a while.
