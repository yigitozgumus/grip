.PHONY: build install clean test run

BINARY=grip
BUILD_DIR=bin

build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY) ./cmd/grip

install: build
	cp $(BUILD_DIR)/$(BINARY) $(HOME)/.local/bin/$(BINARY)

clean:
	rm -rf $(BUILD_DIR)

test:
	go test ./...

run:
	go run ./cmd/grip

# Development helpers
tidy:
	go mod tidy

fmt:
	go fmt ./...

lint:
	golangci-lint run
