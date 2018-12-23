all: build

build:
	GOOS=windows go build -o ngo_win64.exe
	GOOS=linux go build -o ngo_linux64
	GOOS=darwin go build -o ngo_mac64
