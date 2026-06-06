BINARY_NAME=notes
MAIN_PACKAGE_PATH=./cmd/app

.PHONY: tidy
tidy:
	go fmt ./...
	go vet ./...
	go mod tidy

.PHONY: run
run: tidy
	go run ${MAIN_PACKAGE_PATH}/main.go
