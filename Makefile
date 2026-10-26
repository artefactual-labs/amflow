GOPATH=$(shell go env GOPATH)
GOBIN=$(GOPATH)/bin
GOCMD=go
export PATH:=$(GOBIN):$(PATH)

# We don't need make's built-in rules.
MAKEFLAGS += --no-builtin-rules
.SUFFIXES:

.PHONY: help
help:
	@echo "make targets"
	@echo
	@echo "  deps"
	@echo "  check         check all the things"
	@echo "  help          this help message"

.PHONY: deps
deps:
	@echo "Downloading modules..."
	go mod download

.PHONY: clean
clean:
	git clean -f -d -x

.PHONY: build
build:
	CGO_ENABLED=0 \
	  go build \
	    -o dist/amflow \
	    -trimpath \
	    -ldflags="-s -w -X github.com/artefactual-labs/amflow/internal/version.version=try" \
	      github.com/artefactual-labs/amflow

.PHONY: frontend
frontend:
	npm --prefix web ci
	npm --prefix web run build

.PHONY: build
try: build
	dist/amflow version

.PHONY: test
test:
	go test ./...

.PHONY: testrace
testrace:
	go test -race ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: fmt
fmt:
	@echo "gofmt -s -l -d -e ."
	@test -z "$(shell gofmt -s -l -d -e . | tee /dev/stderr)"

.PHONY: staticcheck
staticcheck:
	go run honnef.co/go/tools/cmd/staticcheck@v0.8.2 -checks all,-ST1000,-ST1001 ./...

.PHONY: check
check: test testrace vet fmt staticcheck
