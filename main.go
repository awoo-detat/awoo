package main

import (
	"log"
	"net/http"
	"os"

	"stash.corp.synacor.com/hack/werewolf/handler"
)

func main() {
	handler := handler.New()
	http.HandleFunc("/awoo", handler.Awoo)
	http.HandleFunc("/", handler.Connect)

	file, err := os.OpenFile("werewolf.log", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("running...")
	log.SetOutput(file)
	log.Fatal(http.ListenAndServe(":42300", nil))
}
