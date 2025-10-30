# wof

## wof deprecate

Deprecate one or more Who's On First records.

```
$> ./bin/wof deprecate -h
Deprecate one or more Who's On First records.
Usage:
	 ./bin/wof [options] path(N) path(N)
  -superseded-by-id -1
    	The ID to supersede each record with. If -1 then this flag will be ignored. (default -1)
  -superseded-by-reader-uri -1
    	A valid whosonfirst/go-reader URI used to load records that are doing the superseding. Required if -superseding-id is not -1. (default "null://")
  -superseded-by-writer-uri -1
    	A valid whosonfirst/go-writer URI used to update records that are doing the superseding. Required if -superseding-id is not -1. (default "null://")
```

### Examplea

```
$> wof deprecate 1377455275
```

##### "reader" and "writer" URIs

The `wof` tool has its own internal logic for [deriving paths for reading and writing input documents](https://github.com/whosonfirst/wof-cli?tab=readme-ov-file#paths-and-uris).

That being the case by the time a document URI is resolved at the command layer there is not necessarily enough information to write documents _related_ to the document currently being processed. Further even if that context is known it may not be appropriate. For example, if a document in the `sfomuseum-data-wayfinding` repository is being superseded and the new document is parented by a document in the `sfomuseum-data-architecture` repository (using the `-parent-id` flag) then there is nothing in the `sfomuseum-data-wayfinding` context to know where to find that record.

Hence the `-superseded-byreader-uri` and `-superseded-by-writer-uri` flags. There are expected to be valid [whosonfirst/go-reader.Reader](https://github.com/whosonfirst/go-reader) and [whosonfirst/go-writer.Writer](https://github.com/whosonfirst/go-writer) URIs. For example:

```
$> wof supersede \
	-supereded-by-id 1947304165 \
	-superseded-by-reader-uri repo:///usr/local/data/sfomuseum-data-publicart \
	-superseded-by-writer-uri repo:///usr/local/data/sfomuseum-data-publicart \
	1377455275
```

It's a bit cumbersome but the decision was taken, given the potential for many unrelated moving parts, to be explicit rather than clever.

