help:
	go run main.go --help

test:
	go test -cover ./...

simple:
	go run main.go --destination prefix --bucket s3up-test --cache-control "public, max-age=31536000" --strip 2 `find . -path "./fixtures/*"`

build:
	gox -os="linux darwin windows openbsd" ./...

empty:
	go run main.go
