# go-whosonfirst-iterate-reader

Go package implementing go-whosonfirst-iterate/emitter functionality using `whosonfirst/go-reader.Reader` instances.

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/whosonfirst/go-whosonfirst-iterate-reader.svg)](https://pkg.go.dev/github.com/whosonfirst/go-whosonfirst-iterate-reader/v3)

## Example

Version 3.x of this package introduce major, backward-incompatible changes from earlier releases. That said, migragting from version 2.x to 3.x should be relatively straightforward as a the _basic_ concepts are still the same but (hopefully) simplified. Where version 2.x relied on defining a custom callback for looping over records version 3.x use Go's [iter.Seq2](https://pkg.go.dev/iter) iterator construct to yield records as they are encountered.

For example:

```
import (
	"context"
	"flag"
	"log"

	"github.com/whosonfirst/go-whosonfirst-iterate/v3"
)

func main() {

     	var iterator_uri string

	flag.StringVar(&iterator_uri, "iterator-uri", "reader://?reader=fs%3A%2F%2F%2Fusr%2Flocal%2Fwhosonfirst%2Fgo-whosonfirst-iterate-reader%2Ffixtures%2Fdata". "A registered whosonfirst/go-whosonfirst-iterate/v3.Iterator URI.")
	ctx := context.Background()
	
	iter, _:= iterate.NewIterator(ctx, iterator_uri)
	defer iter.Close()
	
	paths := flag.Args()
	
	for rec, _ := range iter.Iterate(ctx, paths...) {
	    	defer rec.Body.Close()
		log.Printf("Indexing %s\n", rec.Path)
	}
}
```

_Error handling removed for the sake of brevity._

### Version 2.x (the old way)

```
import (
	"context"
	"github.com/whosonfirst/go-whosonfirst-iterate/v2/iterator"
	_ "github.com/whosonfirst/go-whosonfirst-iterate-reader"		
	"io"
	"log"
)

func main() {

	ctx := context.Background()
	
	iter_cb := func(ctx context.Context, path string, fh io.ReadSeeker, args ...interface{}) error {
		log.Println(path)	       
		return nil
	}

	iter_uri := "reader://?reader=fs%3A%2F%2F%2Fusr%2Flocal%2Fwhosonfirst%2Fgo-whosonfirst-iterate-reader%2Ffixtures%2Fdata"
	
	iter, _ := iterator.NewIterator(ctx, iter_uri, emitter_cb)

	iter.IterateURIs(ctx, "102/527/513/102527513.geojson")
``

## Tools

```
$> make cli
go build -mod readonly -ldflags="-s -w" -o bin/count cmd/count/main.go
go build -mod readonly -ldflags="-s -w" -o bin/emit cmd/emit/main.go
```

### count

Count files in one or more whosonfirst/go-whosonfirst-iterate/v3.Iterator sources.

```
$> ./bin/count -h
Count files in one or more whosonfirst/go-whosonfirst-iterate/v3.Iterator sources.
Usage:
	 ./bin/count [options] uri(N) uri(N)
Valid options are:

  -iterator-uri string
    	A valid whosonfirst/go-whosonfirst-iterate/v3.Iterator URI. Supported iterator URI schemes are: cwd://,directory://,featurecollection://,file://,filelist://,geojsonl://,null://,reader://,repo:// (default "repo://")
  -verbose
    	Enable verbose (debug) logging.
```

For example:

```
$> ./bin/count \
	-iterator-uri 'reader://?reader=fs:///usr/local/src/go-whosonfirst-iterate-reader/fixtures/data' \
	102527513.geojson
	
2025/06/26 05:59:32 INFO Counted records count=1 time=873.52µs
```

### emit

Emit records in one or more whosonfirst/go-whosonfirst-iterate/v3.Iterator sources as structured data.

```
$> ./bin/emit -h
Emit records in one or more whosonfirst/go-whosonfirst-iterate/v3.Iterator sources as structured data.
Usage:
	 ./bin/emit [options] uri(N) uri(N)
Valid options are:

  -geojson
    	Emit features as a well-formed GeoJSON FeatureCollection record.
  -iterator-uri string
    	A valid whosonfirst/go-whosonfirst-iterate/v3.Iterator URI. Supported iterator URI schemes are: cwd://,directory://,featurecollection://,file://,filelist://,geojsonl://,null://,reader://,repo:// (default "repo://")
  -json
    	Emit features as a well-formed JSON array.
  -null
    	Publish features to /dev/null
  -stdout
    	Publish features to STDOUT. (default true)
  -verbose
    	Enable verbose (debug) logging.
```

For example:

```
$> go run -mod vendor cmd/emit/main.go \
	-iterator-uri 'reader://?reader=fs:///usr/local/data/sfomuseum-data-whosonfirst/data' 1 \

| jq '.["properties"]["wof:name"]'

"Null Island"
```

_ID `1` is [Null Island](https://spelunker.whosonfirst.org/id/1)._

## See also

* https://github.com/whosonfirst/go-whosonfirst-iterate
* https://github.com/whosonfirst/go-reader 