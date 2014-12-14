# Elasticsearch API Generator

This is **experimental code** (and a hack really) to automatically generate
the Go API from the [REST API specification]((https://github.com/elasticsearch/elasticsearch/tree/master/rest-api-spec))
that comes with Elasticsearch.

You need Go 1.4 for this code relies on the `go generate` tool.

Run `go generate -v && gofmt -w .` to generate and format the services.

