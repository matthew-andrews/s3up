help:
	go run main.go --help

test:
	go test -v -race ./...

simple:
	go run main.go --strip 1 `find . -path "./fixtures/*"`
