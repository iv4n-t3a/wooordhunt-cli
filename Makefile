SOURCES := $(find . -type f -name *.go)
TARGET := wh

build: $(TARGET)

install: $(TARGET)
	@install $(TARGET) /bin/

$(TARGET): $(SOURCES)
	go build -o $(TARGET) cmd/wooordhunt-cli/main.go
