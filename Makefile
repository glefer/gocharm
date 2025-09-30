.PHONY: help build test fmt vet coverage clean

## Makefile for gocharm - minimal helper targets

help: ## Show this help
	@awk 'BEGIN {print "Available targets:"} /^[a-zA-Z0-9_-]+:.*##/ {split($$0,a,":"); desc = substr($$0, index($$0,"##")+3); printf "  %-15s %s\n", a[1], desc}' $(MAKEFILE_LIST)

test: ## Run unit tests
	go test ./...

fmt: ## Run gofmt to format the code in place
	gofmt -w .

vet: ## Run go vet
	go vet ./...

coverage: ## Run tests and write coverage.out
	go test ./... -coverprofile=coverage.out

clean: ## Remove build artifacts
	rm -rf bin coverage.out
