LDFLAGS ?= -s -w
OUTPUT ?= output

.PHONY: all

all: clean build

build:
	mkdir -p "$(OUTPUT)"
	env GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o "$(OUTPUT)"/bput-darwin-arm64
	env GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o "$(OUTPUT)"/bput-linux-amd64
	env GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o "$(OUTPUT)"/bput-windows-amd64.exe

clean:
	rm -f "$(OUTPUT)"/bput-darwin-arm64
	rm -f "$(OUTPUT)"/bput-linux-amd64
	rm -f "$(OUTPUT)"/bput-windows-amd64.exe
