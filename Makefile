all: test assets build

test:
	go test ./...

assets:
	go-bindata-assetfs static/...

build:
	go build
