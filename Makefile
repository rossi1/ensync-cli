.PHONY: build test lint clean run-test test-integration test-all

VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT  ?= $(shell git rev-parse --short HEAD)
DATE    ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -X github.com/ensync-cli/pkg/version.version=$(VERSION) \
           -X github.com/ensync-cli/pkg/version.commit=$(COMMIT) \
           -X github.com/ensync-cli/pkg/version.buildDate=$(DATE)

build:
	go build -ldflags "$(LDFLAGS)" -o bin/ensync

test:
	go test -v -race ./...

test-integration:
	go test -v -race ./test/integration/...

test-all: test test-integration

lint:
	golangci-lint run

clean:
	rm -rf bin/

# Run test commands
run-test: build
	@echo "Testing version command..."
	./bin/ensync version
	
	@echo "\nTesting event commands..."
	@echo '{"key1": "value1", "key2": "value2"}' > test-event.json
	./bin/ensync event create --name "test-event" --payload-file test-event.json
	./bin/ensync event list --page 0 --limit 10 --order DESC
	
	@echo "\nTesting access key commands..."
	@echo '{"send": ["event1"], "receive": ["event3"], "access": ["service1"]}' > test-permissions.json
	./bin/ensync access-key create --file test-permissions.json
	./bin/ensync access-key list
	
	@echo "\nTesting debug mode..."
	ENSYNC_API_KEY="test-key" ./bin/ensync --debug event list
	
	@echo "\nCleaning up test files..."
	rm -f test-event.json test-permissions.json

# Setup test environment
setup-test:
	@echo "Setting up test environment..."
	mkdir -p ~/.ensync
	@echo "base_url: \"http://localhost:8080/api/v1/ensync\"" > ~/.ensync/config.yaml
	@echo "api_key: \"your-test-api-key\"" >> ~/.ensync/config.yaml
	@echo "debug: false" >> ~/.ensync/config.yaml