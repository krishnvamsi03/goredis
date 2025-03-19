BINARY_NAME=goredis
COVERAGE_DIR=coverage

all: clean build

build:
	@echo "building $(BINARY_NAME)"
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -f" -o $(BINARY_NAME) cmd/server/main.go

build-windows:
	@echo "building $(BINARY_NAME) for windows"
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME).exe cmd/server/main.go

clean:
	@echo "cleaning"
	@rm -rf $(COVERAGE_DIR)
	@rm -f $(BINARY_NAME)
	@go clean -i ./...

run:
	@echo "running $(BINARY_NAME)"
	@go run cmd/server/main.go

test:
	@echo "running tests"
	@go test -v -count=1 ./...

coverage:
	@echo "generating coverage report..."
	@mkdir -p $(COVERAGE_DIR)
	@go test -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	@go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "coverage report generated in $(COVERAGE_DIR)/coverage.html"

lint:
	@echo "running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "instanlling golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run ./...; \
	fi

fmt:
	@echo "formatting code"
	@go fmt ./...
	@if command -v goimports > /dev/null; then \
		goimports -w .;\
	else \
		echo "installing goimports..."; \
		go install golang.org/x/tools/cmd/goimports@latest; \
		goimports -w .; \
	fi

check: fmt lint

build-cli:
	@echo "building goredis client"
	@CGO_ENABLED=0 go build -o goredis cmd/cli/main.go

run-cli:
	@go run cmd/cli/main.go
	