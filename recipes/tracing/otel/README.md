# Tracing with OpenTelemetry

Run Jaeger (any other supported tracers work of course):

```sh
./run-tracer.sh
```

Then open the web UI:

```sh
open http://localhost:16686
```

Then run the example, e.g.:

```sh
go build
./tracing -index=test -type=doc -n=100000 -bulk-size=100
```
