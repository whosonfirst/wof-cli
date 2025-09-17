# wof cli cookbook

## Fetching all the `whosonfirst-data-admin-*` repositories as GeoParquet files

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

* It should be possible for the `wof emit` command to produce GeoParquet files natively but this example uses GDAL's `ogr2ogr` to convert files containing line-separated GeoJSON.
