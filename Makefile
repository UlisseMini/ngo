all: clean depends release

clean:
	@rm -rf bin

depends:
	@echo Downloading dependencies...
	GO111MODULE=on go mod download
	@echo Done.

release:
	@mkdir bin
	@echo Building release binaries
	CGO_ENABLED=0 GOOS=windows go build -o bin/ngo_win64.exe
	CGO_ENABLED=0 GOOS=linux go build -o bin/ngo_linux64
	CGO_ENABLED=0 GOOS=darwin go build -o bin/ngo_mac64
	@echo Done.

test:
	go test -cover ./...
