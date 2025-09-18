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

```
#!/bin/sh

touch wof.geojsonl

for REPO in `./bin/wof repos -prefix 'whosonfirst-data-admin-'`
do

    echo "Fetch ${REPO}"

    if [ -f "${REPO}.parquet" ]
    then
	echo "Export exists for ${REPO}"
	continue
    fi
    
    echo "./bin/wof emit -format geojson -iterator-uri git:///tmp https://github.com/whosonfirst-data/${REPO}.git > ${REPO}.geojsonl"
    ./bin/wof emit -format geojson -iterator-uri git:///tmp https://github.com/whosonfirst-data/${REPO}.git >> wof.geojsonl
done

echo "ogr2ogr -f Parquet wof.parquet -oo OGRGeoJSONAllowNonStandard=YES 'GeoJSONSeq:wof.geojsonl'"
ogr2ogr -f Parquet wof.parquet -oo OGRGeoJSONAllowNonStandard=YES "GeoJSONSeq:wof.geojsonl"
```

_Note that this will take "some number of" hours depending on your hardware and network configurations and produce a 22GB line-separated GeoJSON file and an 8GB GeoParquet file._

## Fetching all the `whosonfirst-data-admin-*` repositories as GeoParquet files (native)

_This method has bugs. It will be updated when those issues have been resolved._
