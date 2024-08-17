GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

# Experimental: Feature-based build tags to allow people to build trimmed-down binaries if they want to
TAGS=centroid,emit,emit_geoparquet,export,format,geometry,pip,pip_pmtiles,pip_sqlite,property,show,supersede,uri,validate

cli:
	go build \
		-mod $(GOMOD) \
		-ldflags="$(LDFLAGS)" \
		-tags $(TAGS) \
		-o bin/wof \
		cmd/wof/main.go
