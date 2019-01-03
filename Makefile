.DEFAULT_GOAL := help
.PHONY: test install-linters test-386 test-amd64 lint
 
build-grammar: ## Generate source code for REXT grammar
	nex -s src/rxt/grammar/lexer.nex

mocks: ## Create all mock files for unit tests
	echo "Generating mock files"
	rm -rf ./src/config/mocks/*.go
	mockery -all -dir ./src/config -output ./src/config/mocks
	grep -rl "github.com/denisacostaq/rextporter/src/config" src/config/mocks | xargs sed -i 's:"github.com/denisacostaq/rextporter/src/config":"github.com/simelo/rextporter/src/config":g'

test-grammar: build-grammar ## Test cases for REXT lexer and parser
	go run cmd/rxtc/lexer.go < src/rxt/testdata/skyexample.rxt 2> src/rxt/testdata/skyexample.golden.orig
	diff -u src/rxt/testdata/skyexample.golden src/rxt/testdata/skyexample.golden.orig

test: mocks ## Run test with GOARCH=Default
	go test -count=1 github.com/simelo/rextporter/src/config
	go test -count=1 github.com/simelo/rextporter/src/scrapper
	go test -count=1 github.com/simelo/rextporter/src/memconfig
	if ! screen -list | grep -q "fakeSkycoinForIntegrationTest"; then echo "creating screen fakeSkycoinForIntegrationTest"; screen -L -dm -S fakeSkycoinForIntegrationTest go run test/integration/fake_skycoin_node.go; else echo "fakeSkycoinForIntegrationTest screen already exist. quiting it to create a new one"; screen -S fakeSkycoinForIntegrationTest -X quit; screen -dm -S fakeSkycoinForIntegrationTest go run test/integration/fake_skycoin_node.go; fi
	sleep 3
	go test -count=1 -cpu=1 -parallel=1 github.com/simelo/rextporter/test/integration -args -test.v
	# screen -list can return a not 0 value, this is interpreted as a fail for travis, so use || true
	screen -list || true
	screen -S fakeSkycoinForIntegrationTest -X quit
	cat screenlog.0


test-386: mocks ## Run tests  with GOARCH=386
	GOARCH=386 go test -count=1 github.com/simelo/rextporter/src/config
	GOARCH=386 go test -count=1 github.com/simelo/rextporter/src/scrapper
	GOARCH=386 go test -count=1 github.com/simelo/rextporter/src/memconfig
	if ! screen -list | grep -q "fakeSkycoinForIntegrationTest"; then echo "creating screen fakeSkycoinForIntegrationTest"; screen -L -dm -S fakeSkycoinForIntegrationTest go run test/integration/fake_skycoin_node.go; else echo "fakeSkycoinForIntegrationTest screen already exist. quiting it to create a new one"; screen -S fakeSkycoinForIntegrationTest -X quit; screen -dm -S fakeSkycoinForIntegrationTest go run test/integration/fake_skycoin_node.go; fi
	sleep 3
	GOARCH=386 go test -cpu=1 -parallel=1  -count=1 github.com/simelo/rextporter/test/integration -args -test.v
	# screen -list can return a not 0 value, this is interpreted as a fail for travis, so use || true
	screen -list || true
	screen -S fakeSkycoinForIntegrationTest -X quit
	cat screenlog.0

test-amd64: mocks ## Run tests with GOARCH=amd64
	GOARCH=amd64 go test -count=1 github.com/simelo/rextporter/src/config
	GOARCH=amd64 go test -count=1 github.com/simelo/rextporter/src/scrapper
	GOARCH=amd64 go test -count=1 github.com/simelo/rextporter/src/memconfig
	if ! screen -list | grep -q "fakeSkycoinForIntegrationTest"; then echo "creating screen fakeSkycoinForIntegrationTest"; screen -L -dm -S fakeSkycoinForIntegrationTest go run test/integration/fake_skycoin_node.go; else echo "fakeSkycoinForIntegrationTest screen already exist. quiting it to create a new one"; screen -S fakeSkycoinForIntegrationTest -X quit; screen -dm -S fakeSkycoinForIntegrationTest go run test/integration/fake_skycoin_node.go; fi
	sleep 3
	GOARCH=amd64 go test -cpu=1 -parallel=1  -count=1 github.com/simelo/rextporter/test/integration -args -test.v
	# screen -list can return a not 0 value, this is interpreted as a fail for travis, so use || true
	screen -list || true
	screen -S fakeSkycoinForIntegrationTest -X quit
	cat screenlog.0

lint: ## Run linters. Use make install-linters first.
	vendorcheck ./...
	golangci-lint run -c .golangci.yml ./...
	# The govet version in golangci-lint is out of date and has spurious warnings, run it separately
	go vet -all ./...

check:
	test

install-linters: ## Install linters
	go get -u github.com/FiloSottile/vendorcheck
	# For some reason this install method is not recommended, see https://github.com/golangci/golangci-lint#install
	# However, they suggest `curl ... | bash` which we should not do
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
