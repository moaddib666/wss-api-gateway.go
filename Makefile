# Variables
APP_NAME = margay
CMD_PACKAGE = .

# Build targets
.PHONY: build
build:
	@go build -o $(APP_NAME) $(CMD_PACKAGE)

.PHONY: clean
clean:
	@rm -f $(APP_NAME)

# Test targets
.PHONY: test
test:
	@go test -v ./...

.PHONY: test-coverage
test-coverage:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

# Other targets
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build        - Builds the $(APP_NAME) binary"
	@echo "  clean        - Removes the $(APP_NAME) binary"
	@echo "  test         - Runs all tests in verbose mode"
	@echo "  test-coverage - Runs all tests and generates a coverage report"
	@echo "  help         - Shows this help message"
