CWD=$(shell pwd)

GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

debug:
	go run -mod $(GOMOD) cmd/server/main.go \
		-server-uri http://localhost:8080 \
		-spelunker-uri null://

lambda-server:
	if test -f bootstrap; then rm -f bootstrap; fi
	if test -f server.zip; then rm -f server.zip; fi
	GOARCH=arm64 GOOS=linux go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -tags lambda.norpc -o bootstrap cmd/server/main.go
	zip server.zip bootstrap
	rm -f bootstrap
