default: clean gen

clean:
	rm *.gen.go

gen:
	go generate -v && gofmt -w .
