# wof

## edit

Launch a local web server running on a random port hosting a web application for editing one or more Who's On First documents.

```
$> ./bin/wof edit -h
Launch a local web server running on a random port hosting a web application for editing one or more Who's On First records.
Usage:
	 ./bin/wof path(N) path(N)
```

The `wof edit` command is modeled after the [iandees/wof-editor](https://github.com/iandees/wof-editor) package with the following changes:

* It uses a local WebAssembly (WASM) binary to expose methods for data formatting, data validation and Who's On First placetype-related functionality.

* It provides a user-interface for editing the raw GeoJSON of a Who's On First record.

* It does NOT submit changes as pull requests against a [whosonfirst-data](#) repository. Changes are only written to local file that a record's data was originally read from (or to STDOUT if the record was read from STDIN).

### Examples

```
$> ./bin/wof edit \
	/usr/local/data/whosonfirst/whosonfirst-data-admin-us/data/102/527/513/102527513.geojson \
	/usr/local/data/whosonfirst/whosonfirst-data-admin-us/data/859/225/83/85922583.geojson

2025/10/07 11:10:17 INFO Server is ready and features are viewable url=http://localhost:62963
```

This will open `http://localhost:62963` (or whatever port the tool chooses) in your default web browser to a "list" view, like this:

![](docs/wof-cli-edit-launch.png)

Clicking on one of the links will display that WOF record in the "form" view for editing, like this:

![](docs/wof-cli-edit-form.png)

_You can return to the "list" view by clicking on `WOF Editor` header._

The form view is a limited set of common WOF properties which most often need to be edited. This view is modeled after the interface in the [iandees/wof-editor](https://github.com/iandees/wof-editor) package. If you want to edit the raw GeoJSON for the record itself click the `Data view` button in the navigation tool bar. This will display the record in an editable `<pre>` element like this:

![](docs/wof-cli-edit-data.png)

_The `Format` and `Validate` buttons are only enabled in "data" view. Data validation and formatting happens automatically in "form" view whenever a property is updated._

## Notes (and caveats)

Changes made in the `wof edit` tool are not written to disk (or STDOUT) until the `Save` button is pressed.

It is not possible to edit geometries yet. While there is nothing to prevent you from editing geometries in the raw "data" view those changes will not be reflected until the document being edited is saved and reloaded.

There is a noticeable delay at startup while the WASM binary (used for data formatting, validation and other functionality) is initialized. This is not ideal and future work will focus on speeding it up. For the time being it is an accepted inconvenience.

Currently the tool is hard-coded to use (base) map tiles from OpenStreetMap. Future releases will support map tiles from the Protomaps API or a local Protomaps database file.

## WASM

The `wof edit` tool is bundled with a pre-built WASM binary for performing edit-related functions. If you need or want to rebuild the binary the easiest way to do that is to run the `wasmjs` Makefile target from the root of the `wof-cli` project. For example:

```
$> make wasmjs
GOOS=js GOARCH=wasm \
		go build -mod vendor -ldflags="-s -w" -tags wasmjs \
		-o edit/static/wasm/wof_edit.wasm \
		cmd/wof-edit-wasm/main.go
```

## See also

* https://github.com/iandees/wof-editor
* https://github.com/yairEO/tagify