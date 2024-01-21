[private]
default:
    @just --list --unsorted

# format the codebase
format:
    gofumpt.exe -l -w .
    golines.exe -l -w .

# run linters on the codebase
lint:
    golangci-lint run ./...

# run all unit tests
test:
    go test -v -cover -coverprofile coverage.out ./...

# build the application binary
build:
    go build -v ./...
