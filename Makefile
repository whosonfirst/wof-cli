GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

cli:
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/wof cmd/wof/main.go

# https://github.com/marcboeker/go-duckdb?tab=readme-ov-file#vendoring
modvendor:
        modvendor -copy="**/*.a **/*.h" -v
