help:
	go run main.go --help

test:
	go test -v -race ./...

simple:
	go run main.go --bucket s3up-test --cache-control "public, max-age=31536000" --strip 1 `find . -path "./fixtures/one-file/*"`

empty:
	go run main.go
