BUILD_TIMESTAMP=$(shell date -Iseconds)
BUILD_COMMIT=$(shell git rev-parse --short HEAD)
BUILD_VERSION=latest
BINARY=routemap
LDFLAGS=-ldflags "-X main.version=$(BUILD_VERSION) -X main.commit=$(BUILD_COMMIT) -X main.date=$(BUILD_TIMESTAMP)"

all: build

.PHONY: test
test:
	go test -v ./...

# Builds for the local platform only.
.PHONY: build
build:
	go build $(LDFLAGS) -o build/$(BINARY) ./cmd/routemap

# Run the goreleaser process to produce a snapshot build (no publishing).
.PHONY: snapshot
snapshot:
	goreleaser --skip-publish --snapshot

# Cleans up local platform and goreleaser builds.
.PHONY: clean
clean:
	rm -rf build/
	rm -rf dist/

