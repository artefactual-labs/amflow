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

TOOLS= \
	github.com/goadesign/goa/goagen \
	github.com/gobuffalo/packr/v2/packr2

.PHONY: deps
deps:
	@echo "Downloading modules..."
	go mod download
	@echo "Installing tools..."
	CGO_ENABLED=0 go install $(TOOLS)

.PHONY: goagen
goagen:
	@goagen app     -d github.com/sevein/amflow/design -o internal/api
	@goagen swagger -d github.com/sevein/amflow/design -o public
	@goagen schema  -d github.com/sevein/amflow/design -o public
	@goagen js      -d github.com/sevein/amflow/design -o ui/js/client --noexample

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
	    -ldflags="-s -X github.com/sevein/amflow/internal/version.version=try" \
	      github.com/sevein/amflow

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
