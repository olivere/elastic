default: clean gen

clean:
	rm -f *.gen.go

gen:
	go generate -v && gofmt -w .
