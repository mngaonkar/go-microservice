package main

import (
	"embed"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/urfave/negroni"
)

//go:embed static/*
var content embed.FS

type RecommendRequest struct {
	Name string `json:"name"`
}

// send single page application to client
func index(w http.ResponseWriter, r *http.Request) {
	log.Println("request for index ", *r)
	rawFile, err := content.ReadFile("static/index.html")
	if err != nil {
		log.Fatal("error reading index.html, err = ", err)
		return
	}
	w.Write(rawFile)
}

// send recommendation
func recommend(w http.ResponseWriter, r *http.Request) {
	log.Println("recommend request received, r = ", *r)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("error reading POST payload, err = ", err)
		return
	}

	var request RecommendRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		log.Fatal("error parsing POST payload, err = ", err)
		return
	}
	err = postRecommendRequest(request)
	if err != nil {
		log.Fatal("error sending request to Kafka")
	}
}

// post request to kafka
func postRecommendRequest(request RecommendRequest) error {
	log.Printf("received recommendation like %s", request.Name)
	return nil
}

func main() {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	mux.HandleFunc("/", index)
	mux.HandleFunc("/v1/recommender", recommend)

	n := negroni.Classic()
	n.UseHandler(mux)
	http.ListenAndServe(":3000", n)
}
