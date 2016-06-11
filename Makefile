clean:
	git clean -fxd

help:
	go run main.go --help

test:
	go test -cover ./...

simple:
	find . -path "./fixtures/*" -exec go run main.go --destination prefix --bucket s3up-test --cache-control "public, max-age=31536000" --strip 2 {} +

dry:
	find . -path "./fixtures/*" -exec go run main.go --destination prefix --bucket s3up-test --cache-control "public, max-age=31536000" --strip 2 --dry-run {} +

build:
	gox -os="linux darwin windows openbsd" ./...

empty:
	go run main.go
