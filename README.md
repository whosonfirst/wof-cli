# wof-cli

Experimental standalone 'wof' command-line tool for common Who's On First operations.

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
* pip
* validate
```

_Important: The inputs and outputs for the `wof` tool have not been finalized yet, notably about how files are read and written if updated. You should expect change in the short-term._

#### wof export

#### wof fmt

#### wof pip

#### wof validate

## See also

* https://github.com/whosonfirst/go-whosonfirst-export
* https://github.com/whosonfirst/go-whosonfirst-format
* https://github.com/whosonfirst/go-whosonfirst-validate
* https://github.com/whosonfirst/go-whosonfirst-spatial