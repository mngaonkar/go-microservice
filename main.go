package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/urfave/negroni"
)

//go:embed all:static
var StaticFiles embed.FS

type RecommendRequest struct {
	name string
}

// send single page application to client
func index(w http.ResponseWriter, r *http.Request) {
	rawFile, _ := StaticFiles.ReadFile("index.html")
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
