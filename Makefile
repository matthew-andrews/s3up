help:
	go run main.go --help

test:
	go test -v -race ./...

simple:
	go run main.go --destination prefix --bucket s3up-test --cache-control "public, max-age=31536000" --strip 2 `find . -path "./fixtures/one-file/*"`

empty:
	go run main.go
