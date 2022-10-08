SHELL=/usr/bin/env bash

# Project specific properties.
application_name           = ledgerconv
application_binary_name    = ledgerconv

# For ProtoBuf code generation.
proto_path=src/proto/*.proto

# Builds the project.
build:
	@echo "+$@"
	@go build -o bin/$(application_binary_name)

# Runs the project after linting and building it anew.
run: tidy lint build
	@echo "+$@"
	@echo "########### Running the application binary ############"
	@bin/$(application_binary_name)

# Tests the whole project.
test:
	@echo "+$@"
	@CGO_ENABLED=1 go test -race -coverprofile=coverage.out -covermode=atomic ./...

# Runs the "go mod tidy" command.
tidy:
	@echo "+$@"
	@go mod tidy

# Runs golang-ci-lint over the project.
lint:
	@echo "+$@"
	@golangci-lint run

# Generates code using the found protocol buffer files.
proto:
	@echo "+$@"
	@protoc \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		$(proto_path)
