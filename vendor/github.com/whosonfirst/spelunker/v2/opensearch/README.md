# opensearch

The `opensearch` package implements the `Spelunker` interface for Who's On First data indexed in an [OpenSearch](https://opensearch.org/) database.

## Connecting to OpenSearch

Ultimately, all of the code in this package uses the `opensearch-project/opensearch-go/v4` package for executing requests against an OpenSearch server. It also uses the `whosonfirst/go-whosonfirst-database/opensearch/client` package for managing the details of creating an `opensearch-go` client instance. These clients are derived from URIs which take the form of:

```
opensearch://{OPENSEARCH_HOST}:{OPENSEARCH_PORT}/{OPENSEARCH_INDEX}?{QUERY_PARAMETERS}
```

Where {QUERY_PARAMETERS} may be one or more of the following:
* `debug={BOOLEAN}`. A boolean value to configure the underlying OpenSearch client to write request and response bodies to STDOUT.
* `insecure={BOOLEAN}`. A boolean value to disable TLS "InsecureSkipVerify" checks (for custom certificate authorities and the like).
* `require-tls={BOOLEAN}`. A boolean value to ensure that all connections are made over HTTPS even if the OpenSearch port is not 443.
* `username={STRING}`. The OpenSearch username for authenticated connections.
* `password={STRING}`. The OpenSearch password for authenticated connections.
* `aws-credentials-uri={STRING}`. A a valid `aaronland/go-aws-auth` URI used to create a Golang AWS authentication config used to sign requests to an AWS-hosted OpenSearch instance.
* `bulk-index={BOOLEAN}`. A boolean value. If true then writes will be performed using a "bulk indexer". Default is true.
* `workers={INT}`. The number of users to enable for bulk indexing. Default is 10.

These are the common query parameters when creating clients for both the `wof-spelunker-httpd` and `wof-spelunker-index` tools. The `wof-spelunker-index` tool also accepts the following parameters:

* `bulk-index={BOOLEAN}`. A boolean value. If true then writes will be performed using a "bulk indexer". Default is true.
* `workers={INT}`. The number of users to enable for bulk indexing. Default is 10.

## Examples

### Running locally

These examples assume a "local" setup meaning there is local instance of OpenSearch running on port 9200. The "easiest" way to do this is with the [Docker](https://www.docker.com/) application running a containerized instance of OpenSearch and the `os-local` Makefile target provided by this package.

For example, in one terminal window:

```
$> cd spelunker
$> make os-local
docker run \
		-it \
		-p 9200:9200 \
		-p 9600:9600 \
		-e "discovery.type=single-node" \
		-e "OPENSEARCH_INITIAL_ADMIN_PASSWORD=dkjfhsjdkfkjdjhksfhskd98475kjHkzjxckj" \
		-v opensearch-data1:/usr/local/data/opensearch \
		opensearchproject/opensearch:latest

...wait for Docker/OpenSearch to start
```

In another terminal run the `os-local-index` Makefile target passing in one or more Who's On First data sources in the `REPOS` variable. This will index the records in those sources in to the OpenSearch instance started, above. For example:

```
$> make os-local-index REPOS=/usr/local/data/whosonfirst/whosonfirst-data-admin-ca
go run -tags opensearch -mod readonly ./cmd/wof-spelunker-index/main.go opensearch \
		-create-index \
		-client-uri 'opensearch2://localhost:9200/spelunker?username=admin&password=dkjfhsjdkfkjdjhksfhskd98475kjHkzjxckj&insecure=true&require-tls=true' \
		/usr/local/data/whosonfirst/whosonfirst-data-admin-ca

2025/11/15 17:34:59 INFO Iterator stats elapsed=17.295467917s seen=33845 allocated="229 MB" "total allocated"="15 GB" sys="643 MB" numgc=182
2025/11/15 17:35:08 INFO Index complete indexed=28097
```

_Note the default OpenSearch username and password values. This are assigned by default (and can be overriden) as Makefile variables._

Once complete run the `os-server-local` Makefile target to start the Spelunker web application reading data from the Spelunker OpenSearch instance:

```
$> make os-local-server
go run -tags opensearch -mod vendor ./cmd/wof-spelunker-httpd/main.go \
		-server-uri http://localhost:8080 \
		-spelunker-uri 'opensearch://?client-uri=https%3A%2F%2Flocalhost%3A9200%2Fspelunker%3Fusername%3Dadmin%26password%3Ddkjfhsjdkfkjdjhksfhskd98475kjHkzjxckj%26insecure%3Dtrue%26require-tls%3Dtrue&cache-uri=ristretto%3A%2F%2F&reader-uri=https%3A%2F%2Fdata.whosonfirst.org'

2025/11/15 11:42:44 INFO Listening for requests address=http://localhost:8080
```

_Note: The value of the `-spelunker-uri` flag is NOT the same as the "client-uri" URI used to connect to OpenSearch. Specifically, the "client-uri" URI is encoded as a query parameter of the `-spelunker-uri` flag. There are additional OpenSearch-implementation-specific flags (`cache-uri` and `reader-uri`). Consult the [cmd/wof-spelunker-httpd documentation](../cmd/wof-spelunker-httpd) for details._

## Things the `opensearch` Spelunker implementation does NOT do yet

* The `opensearch` Spelunker does not implement any of the tag-related methods (`GetTags`, `HasTag`, `HasTagFaceted`) yet.

## Database schema(s)

Database schemas (mappings) used by the `OpenSearchSpelunker` implementation are defined in the [whosonfirst/go-whosonfirst-database/opensearch/schema/v2](https://github.com/whosonfirst/go-whosonfirst-database/tree/main/opensearch/schema/v2) package.