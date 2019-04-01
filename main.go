package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/awoo-detat/awoo/handler"
)

var LocalAssetDir string

func main() {
	handler := handler.New()
	http.HandleFunc("/awoo", handler.Connect)

	var server http.FileSystem
	if LocalAssetDir != "" {
		server = http.Dir(LocalAssetDir)
	} else {
		server = assetFS()
	}
	http.Handle("/", http.FileServer(server))

	port := os.Getenv("PORT")
	if port == "" {
		port = "42300"
	}
	log.Printf("running on port %s...\n", port)
	log.Fatal(http.ListenAndServe(net.JoinHostPort("", port), nil))
}
