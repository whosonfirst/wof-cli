# wof

## wof supersede

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
  -parent-reader-uri string
    	A valid whosonfirst/go-reader URI used to load parent records if -superseded-id is -1. Required if -parent-id is not `-1`. (default "null://")
  -superseding-id int
    	The ID to supersede each record with. If -1 then each record will be cloned and the new ID of the clone will be used as the superseding_id. (default -1)
  -superseding-reader-uri string
    	A valid whosonfirst/go-reader URI used to load records that are doing the superseding. Required if -superseding-id is not -1. (default "null://")
  -superseding-writer-uri string
    	A valid whosonfirst/go-writer URI used to update records that are doing the superseding. Required if -superseding-id is not -1. (default "null://")
```

### Examples

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

