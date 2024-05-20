# wof

Experimental standalone 'wof' binary for common Who's On First operations.

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
* export
* fmt
* validate
```

_Important: The inputs and outputs for the `wof` tool have not been finalized yet, notably about how files are read and written if updated. You should expect change in the short-term._

## See also

* https://github.com/whosonfirst/go-whosonfirst-export
* https://github.com/whosonfirst/go-whosonfirst-format
* https://github.com/whosonfirst/go-whosonfirst-validate