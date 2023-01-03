package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/urfave/negroni"
)

//go:embed static/*
var content embed.FS

type RecommendRequest struct {
	name string
}

// send single page application to client
func index(w http.ResponseWriter, r *http.Request) {
	log.Println("request for index ", *r)
	rawFile, err := content.ReadFile("index.html")
	if err != nil {
		log.Fatal("error reading index.html")
		return
	}
	w.Write(rawFile)
}

// send recommendation
func recommend(w http.ResponseWriter, r *http.Request) {
	var request RecommendRequest
	request.name = "titanic"
	err := postRecommendRequest(request)
	if err != nil {
		log.Fatal("error sending request to Kafka")
	}
}

// post request to kafka
func postRecommendRequest(request RecommendRequest) error {
	return nil
}

func main() {
	mux := http.NewServeMux()
	var fs fs.FS = os.DirFS("static")
	fileServer := http.FileServer(http.FS(fs))
	mux.Handle("/static", fileServer)
	mux.HandleFunc("/", index)
	mux.HandleFunc("/v1/recommend", recommend)

	n := negroni.Classic()
	n.UseHandler(mux)
	http.ListenAndServe(":3000", n)
}
