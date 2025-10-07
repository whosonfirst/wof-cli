# wof

## edit

Launch a local web server running on a random port hosting a web application for editing one or more Who's On First documents. It is modeled after the [iandees/wof-editor](https://github.com/iandees/wof-editor) package with the following changes:

```
$> ./bin/wof edit -h
Launch a local web server running on a random port hosting a web application for editing one or more Who's On First records.
Usage:
	 ./bin/wof path(N) path(N)
```

### Notes

### Examples

```
$> ./bin/wof edit \
	/usr/local/data/whosonfirst/whosonfirst-data-admin-us/data/102/527/513/102527513.geojson \
	/usr/local/data/whosonfirst/whosonfirst-data-admin-us/data/859/225/83/85922583.geojson
```

## See also

* https://github.com/iandees/wof-editor
* https://github.com/yairEO/tagify