.DEFAULT_GOAL := help
.PHONY: test install-linters test-386 test-amd64 lint

test: ## Run test with GOARCH=Default
	# go test  -timeout=5m ./src/...
	# go test  -timeout=5m ./cmd/...
	go test github.com/simelo/rextporter/src/client
	go test github.com/simelo/rextporter/test/integration

test-386: ## Run tests  with GOARCH=386
	# GOARCH=386 go test ./cmd/... -timeout=5m
	# GOARCH=386 go test ./src/... -timeout=5m
	GOARCH=386 go test github.com/simelo/rextporter/src/client
	GOARCH=386 go test github.com/simelo/rextporter/test/integration

test-amd64: ## Run tests with GOARCH=amd64
	# GOARCH=amd64  go test ./cmd/... -timeout=5m
	# GOARCH=amd64  go test ./src/... -timeout=5m
	GOARCH=amd64 go test github.com/simelo/rextporter/src/client
	GOARCH=amd64 go test github.com/simelo/rextporter/test/integration

lint: ## Run linters. Use make install-linters first.
	vendorcheck ./...
	golangci-lint run -c .golangci.yml ./...
	# The govet version in golangci-lint is out of date and has spurious warnings, run it separately
	go vet -all ./...

install-linters: ## Install linters
	go get -u github.com/FiloSottile/vendorcheck
	# For some reason this install method is not recommended, see https://github.com/golangci/golangci-lint#install
	# However, they suggest `curl ... | bash` which we should not do
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
