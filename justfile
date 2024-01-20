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

# build the application binary
build:
    go build

# run all unit tests
test:
    go test -cover -coverprofile coverage.out ./...
