GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

TAGS=sqlite3

vuln:
	govulncheck -show verbose ./...

cli:
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -tags $(TAGS) -o bin/wof cmd/wof/main.go

wasmjs:
	GOOS=js GOARCH=wasm \
		go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -tags wasmjs \
		-o edit/static/wasm/wof_edit.wasm \
		cmd/wof-edit-wasm/main.go

# This compiles but not produce code that executes successfully...
tiny:
	tinygo build -tags wasmjs -target wasm \
		-o edit/static/wasm/wof_edit.wasm \
		cmd/wof-edit-wasm/main.go

# https://github.com/marcboeker/go-duckdb?tab=readme-ov-file#vendoring
modvendor:
	modvendor -copy="**/*.a **/*.h" -v
