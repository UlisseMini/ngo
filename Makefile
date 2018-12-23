all: clean depends build

clean:
	@rm -rf bin

depends:
	GO111MODULE=on go mod download

build:
	@mkdir bin
	CGO_ENABLED=0 GOOS=windows go build -o bin/ngo_win64.exe
	CGO_ENABLED=0 GOOS=linux go build -o bin/ngo_linux64
	CGO_ENABLED=0 GOOS=darwin go build -o bin/ngo_mac64
