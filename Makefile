.PHONY: build clean


build:
	go build -o dist/jum ./cmd/jumie
	go build -o dist/jumied ./cmd/jumied

clean:
	rm -rf dist/
