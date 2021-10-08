# Illustrate Elastic OpenTracing

Run Jaeger (any other supported tracers work of course):

```
$ ./run-tracer.sh
```

Then open the web UI:

```
$ open http://localhost:16686
```

Then run the example, e.g.:

```
$ dep ensure # not necessary for Go 1.11 or later
$ go build
$ ./tracing -index=test -type=doc -n=100000 -bulk-size=100
```
