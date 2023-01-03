package main

import (
	"embed"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/urfave/negroni"
	kafka "microservice.mngaonkar.com/kafka_request_handler"
)

//go:embed static/*
var content embed.FS

var bootstrapServer string = "147.182.230.45:9092"

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
	handler := kafka.NewKafkaWriteHandler(bootstrapServer)
	if handler != nil {
		log.Fatal("failed to create Kafka write handler")
	}

	message, err := json.Marshal(request)
	if err != nil {
		log.Fatal("error marshalling JSON request, err = ", err)
	}

	err = handler.WriteMessage(string(message))
	if err != nil {
		log.Fatal("error writing message to Kafka, err = ", err)
		return err
	}

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
