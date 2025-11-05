CWD=$(shell pwd)

GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

GOTAGS=icu json1 fts5

cli:
	go build -mod $(GOMOD) -tags="$(GOTAGS)" -ldflags="$(LDFLAGS)" -o bin/httpd cmd/httpd/main.go

SPELUNKER_DATABASE=/usr/local/data/ca-search.db
SPELUNKER_URI=sql://sqlite3?dsn=file:$(SPELUNKER_DATABASE)

server:
	go run -mod $(GOMOD) -tags "$(GOTAGS)" cmd/httpd/main.go \
		-server-uri http://localhost:8080 \
		-protomaps-api-key '$(APIKEY)' \
		-spelunker-uri $(SPELUNKER_URI)

		# -spelunker-uri $(SPELUNKER_URI) 

get_id:
	go run -mod $(GOMOD) cmd/spelunker/main.go \
		-spelunker-uri $(SPELUNKER_URI) \
		-command id \
		-id $(ID)

get_descendants:
	go run -mod $(GOMOD) cmd/spelunker/main.go \
		-spelunker-uri $(SPELUNKER_URI) \
		-command descendants \
		-id $(ID)

search:
	go run -mod $(GOMOD) -tags "icu json1 fts5" cmd/spelunker/main.go \
		-spelunker-uri $(SPELUNKER_URI) \
		-command search \
		-query $(NAME)
