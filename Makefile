SOURCES := $(shell find . -type f -name *.go)
TARGET := wh

build: $(TARGET)

install: $(TARGET)
	@install $(TARGET) /bin/

test:
	go mod tidy
	go test -v ./test/...

$(TARGET): $(SOURCES)
	go mod tidy
	go build -o $(TARGET) cmd/wooordhunt-cli/main.go

.PHONY: build, install, test
