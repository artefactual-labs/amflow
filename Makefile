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
	@echo "  goagen"
	@echo "  check         check all the things"
	@echo "  help          this help message"

tools:
	env GO111MODULE=off go get github.com/gobuffalo/packr/packr2

.PHONY: deps
deps: tools
	@echo "Downloading modules..."
	go mod download

.PHONY: goagen
goagen:
	@goagen app     -d github.com/artefactual-labs/amflow/design -o internal/api
	@goagen swagger -d github.com/artefactual-labs/amflow/design -o public
	@goagen schema  -d github.com/artefactual-labs/amflow/design -o public
	@goagen js      -d github.com/artefactual-labs/amflow/design -o web/js/client --noexample

.PHONY: clean
clean:
	git clean -f -d -x

.PHONY: prebuild
prebuild: frontend generate

.PHONY: build
build: prebuild
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOFLAGS=-ldflags=-w \
	  go build \
	    -o dist/amflow \
	    -ldflags="-s -X github.com/artefactual-labs/amflow/internal/version.version=try" \
	      github.com/artefactual-labs/amflow

.PHONY: generate
generate:
	go generate

.PHONY: frontend
frontend:
	yarn --cwd web install
	yarn --cwd web build

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
	staticcheck -checks all,-ST1000,-ST1001 ./...

.PHONY: check
check: test testrace vet fmt staticcheck
