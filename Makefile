.PHONY: build clean


build:
	go build -o dist/juun ./cmd/juun
	go build -o dist/juund ./cmd/juund

clean:
	rm -rf dist/
