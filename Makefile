.DEFAULT_GOAL := buildq

#
# Boring
#

fmt:
	go fmt ./...

lint: fmt
	golint ./...

vet: fmt
	go vet ./...

bootstrap:
	go get github.com/mitchellh/gox
	go generate -tags tools tools/tools.go

#
# Run, test, & bench
#

run:
	go run .

test:
	go test --cover ./...

bench:
	go test --cover -bench . -benchmem ./...

# Reflex doesn't work on windows :(
# @TODO: implement an equivalent file watcher
#watch:
#	reflex -s -r '*.go' make run

#
# Build
#

buildq:
	go build -ldflags "-s -w" .

build: vet
	gox -osarch "linux/amd64 windows/amd64" \
	-gocmd go                  \
	-ldflags "-s -w"           \
	-tags "go-file-search"       \
	-output "bin/go-file-search" .
