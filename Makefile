.PHONY: build clean


build:
	go build -o dist/jumie ./cmd/jumie
	go build -o dist/jumied ./cmd/jumied

clean:
	rm -rf dist/
