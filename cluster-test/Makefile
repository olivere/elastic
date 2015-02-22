.PHONY: build run-omega-cluster-test

default: build

build:
	go build cluster-test.go

run-omega-cluster-test:
	go run cluster-test.go -nodes=http://192.168.2.65:8200,http://192.168.2.64:8200 -n=5 -sniffer=10 -healthcheck=5
