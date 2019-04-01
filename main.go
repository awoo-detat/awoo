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

	file, err := os.OpenFile("werewolf.log", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "42300"
	}
	log.Printf("running on port %s...\n", port)
	log.SetOutput(file)
	log.Fatal(http.ListenAndServe(net.JoinHostPort("", port), nil))
}
