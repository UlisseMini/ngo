all: clean depends test release

LDFLAGS=-ldflags="-w -s"

clean:
	@rm -rf bin

depends:
	GO111MODULE=on go mod download

release:
	@mkdir bin
	@echo Building release binaries
	CGO_ENABLED=0 GOOS=windows go build $(LDFLAGS) -o bin/ngo_win64.exe
	CGO_ENABLED=0 GOOS=linux go build $(LDFLAGS) -o bin/ngo_linux64
	CGO_ENABLED=0 GOOS=darwin go build $(LDFLAGS) -o bin/ngo_mac64

upx:
	upx --brute --best bin/*

test:
	@go test -race -cover ./...
