package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Sigafoos/awoo/handler"
)

func main() {
	handler := handler.New()
	http.HandleFunc("/awoo", handler.Connect)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Print("servin index?")
		http.ServeFile(w, r, "index.html")
	})

	jsfs := http.FileServer(http.Dir("js"))
	http.Handle("/js/", http.StripPrefix("/js", jsfs))
	assetsfs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets", assetsfs))

	file, err := os.OpenFile("werewolf.log", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("running...")
	log.SetOutput(file)
	log.Fatal(http.ListenAndServe(":42300", nil))
}
