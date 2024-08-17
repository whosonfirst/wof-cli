# wof

## wof property

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

### Examples

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
