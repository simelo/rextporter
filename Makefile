.DEFAULT_GOAL := help
.PHONY: test install-linters test-386 test-amd64 lint

test: ## Run test with GOARCH=Default
	# go test  -timeout=5m ./src/...
	# go test  -timeout=5m ./cmd/...
	go test -count=1 github.com/simelo/rextporter/test/config
	go test -count=1 github.com/simelo/rextporter/src/client
	go test -count=1 github.com/simelo/rextporter/test/integration
	@mkdir -p coverage/
	go test -count=1 github.com/simelo/rextporter/test/config
	go test -count=1 github.com/simelo/rextporter/src/client
	screen -dm -S fakeSkycoinForIntegrationTest go run test/integration/fake_skycoin_node.go
	sleep 3
	go test -count=1 -cpu=1 -parallel=1 github.com/simelo/rextporter/test/integration
	screen -S fakeSkycoinForIntegrationTest -X quit

test-386: ## Run tests  with GOARCH=386
	# GOARCH=386 go test ./cmd/... -timeout=5m
	# GOARCH=386 go test ./src/... -timeout=5m
	GOARCH=386 go test -count=1 github.com/simelo/rextporter/src/config
	GOARCH=386 go test -count=1 github.com/simelo/rextporter/src/client
	screen -dm -S fakeSkycoinForIntegrationTest go run test/integration/fake_skycoin_node.go
	sleep 3
	GOARCH=386 go test -cpu=1 -parallel=1  -count=1 github.com/simelo/rextporter/test/integration
	screen -S fakeSkycoinForIntegrationTest -X quit

test-amd64: ## Run tests with GOARCH=amd64
	# GOARCH=amd64  go test ./cmd/... -timeout=5m
	# GOARCH=amd64  go test ./src/... -timeout=5m
	GOARCH=amd64 go test -count=1 github.com/simelo/rextporter/src/config
	GOARCH=amd64 go test -count=1 github.com/simelo/rextporter/src/client
	screen -dm -S fakeSkycoinForIntegrationTest go run test/integration/fake_skycoin_node.go
	sleep 3
	GOARCH=amd64 go test -cpu=1 -parallel=1  -count=1 github.com/simelo/rextporter/test/integration
	screen -S fakeSkycoinForIntegrationTest -X quit

lint: ## Run linters. Use make install-linters first.
	vendorcheck ./...
	golangci-lint run -c .golangci.yml ./...
	# The govet version in golangci-lint is out of date and has spurious warnings, run it separately
	go vet -all ./...

check:
	test integration-test-enable-seed-api

integration-test-enable-seed-api: ## Run enable seed api integration test
	GOCACHE=off COIN=$(COIN) ./ci-scripts/integration-test-enable-seed-api.sh

install-linters: ## Install linters
	go get -u github.com/FiloSottile/vendorcheck
	# For some reason this install method is not recommended, see https://github.com/golangci/golangci-lint#install
	# However, they suggest `curl ... | bash` which we should not do
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'