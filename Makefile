SOURCES := $(shell find . -type f -name *.go)
TARGET := wh

build: $(TARGET)

install: $(TARGET)
	@mkdir -p ~/.config/wh/
	@install config/config.json ~/.config/wh/
	@install $(TARGET) ~/.local/bin/

test:
	go mod tidy
	go test -v ./test/...

$(TARGET): $(SOURCES)
	go mod tidy
	go build -o $(TARGET) cmd/wooordhunt-cli/main.go

.PHONY: build, install, test
