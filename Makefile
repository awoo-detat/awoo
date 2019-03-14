all: test assets build

test:
	go test ./...

elm:
	elm-make src/Main.elm --output static/js/awoo.js

assets: elm
	go-bindata-assetfs static/...

run:
	go run -ldflags '-X main.LocalAssetDir=static' bindata.go main.go

build:
	go build
