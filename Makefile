.PHONY: all
all: clean test

.PHONY: clean
clean:
	go clean
	go mod tidy

.PHONY: test
test:
	go test -v ./...
