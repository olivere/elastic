# Cluster Test

This directory contains a program you can use to test a cluster.

Here's how:

First, install a cluster of Elasticsearch nodes. You can install them on
different computers, or start several nodes on a single machine.

Build cluster-test by `go build cluster-test.go`.

Run `./cluster-test -h` to get a list of flags:

```sh
$ ./cluster-test -h
Usage of ./cluster-test:
  -healthcheck=-1: healthcheck schedule (in seconds)
  -index="twitter": name of ES index to use
  -log="": log file
  -n=5: number of goroutines that run searches
  -nodes="": comma-separated list of ES URLs (e.g. 'http://192.168.2.10:9200,http://192.168.2.11:9200')
  -retries=5: number of retries
  -sniffer=-1: sniffer schedule (in seconds)
  -trace="": trace file
```

Example:

```sh
$ ./cluster-test` with your cluster settings, e.g.
`./cluster-test -nodes=http://127.0.0.1:9200,http://127.0.0.1:9201,http://127.0.0.1:9202 -n=5 -retries=5 -sniffer=10 -healthcheck=5 -log=elastic.log
```

The above example will starts on the cluster defined by the three nodes,
available as http://127.0.0.1:9200, http://127.0.0.1:9201, and
http://127.0.0.1:9202.

It will run 5 search jobs in parallel (`-n=5`).

It will retry failed requests 5 times (`-retries=5`).

It will sniff the cluster every 10 seconds (`-sniffer=10`).

It will perform a health check on the known nodes every 5 seconds (`-healthcheck=5`).

It will write a log file (`-log=elastic.log`).

