# wasm

The `wof-edit.wasm` binary is built using the `wasmjs` Makefile target.

```
$> make wasmjs
GOOS=js GOARCH=wasm \
		go build -mod readonly -ldflags="-s -w" -tags wasmjs \
		-o edit/static/wasm/wof_edit.wasm \
		cmd/wof-edit-wasm/main.go
```
