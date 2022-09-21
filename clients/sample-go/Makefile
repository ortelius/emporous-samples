GO := go

GO_BUILD_PACKAGES := ./cmd/sample-client
GO_BUILD_BINDIR :=./bin

build: prep-build-dir
	$(GO) build -o $(GO_BUILD_BINDIR)/sample-client  -ldflags="$(GO_LD_EXTRAFLAGS)" $(GO_BUILD_PACKAGES)
.PHONY: build

prep-build-dir:
	mkdir -p ${GO_BUILD_BINDIR}
.PHONY: prep-build-dir

vendor:
	$(GO) mod tidy
	$(GO) mod verify
	$(GO) mod vendor
.PHONY: vendor

clean:
	@rm -rf ./$(GO_BUILD_BINDIR)/*
.PHONY: clean

test-unit:
	$(GO) test $(GO_BUILD_FLAGS) -coverprofile=coverage.out -race -count=1 ./...
.PHONY: test-unit

sanity: vendor format vet
	git diff --exit-code
.PHONY: sanity

format: 
	$(GO) fmt ./...
.PHONY: format

vet: 
	$(GO) vet ./...
.PHONY: vet

all: clean vendor test-unit build
.PHONY: all
