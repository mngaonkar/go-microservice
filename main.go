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
var request_topic string = "recommend-requests"
var response_topic string = "recommend-responses"

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

// get recommendation response
func getRecommendResponse(w http.ResponseWriter, r *http.Request) {
	log.Println("recommend get response received, r = ", *r)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("error reading POST payload, err = ", err)
		return
	}

	var request GetRecommendResponse
	err = json.Unmarshal(body, &request)
	if err != nil {
		log.Fatal("error parsing POST payload, err = ", err)
		return
	}

	// send recommendation response request
	if r.Method == http.MethodPost {

	}
}

// send recommendation
func sendRecommendRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("recommend request received, r = ", *r)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("error reading POST payload, err = ", err)
		return
	}

	var request SendRecommendRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		log.Fatal("error parsing POST payload, err = ", err)
		return
	}

	// send recomendation request
	if r.Method == http.MethodPost {
		err = postRecommendRequest(request)
		if err != nil {
			log.Fatal("error sending request to Kafka")
		}
	} else {
		log.Fatal("invalid request")
	}
}

func postRecommendResponse(request GetRecommendResponse) error {
	log.Printf("received get response request for %s", request.Name)
	handler := kafka.NewKafkaWriteHandler(bootstrapServer)
	if handler == nil {
		log.Fatal("failed to create Kafka write handler")
	}

	message, err := json.Marshal(request)
	if err != nil {
		log.Fatal("error marshalling JSON request, err = ", err)
	}

	err = handler.WriteMessage(string(message), response_topic)
	if err != nil {
		log.Fatal("error writing message to Kafka, err = ", err)
		return err
	}

	return nil
}

// post request to kafka
func postRecommendRequest(request SendRecommendRequest) error {
	log.Printf("received recommendation like %s", request.Name)
	handler := kafka.NewKafkaWriteHandler(bootstrapServer)
	if handler == nil {
		log.Fatal("failed to create Kafka write handler")
	}

	message, err := json.Marshal(request)
	if err != nil {
		log.Fatal("error marshalling JSON request, err = ", err)
	}

	err = handler.WriteMessage(string(message), request_topic)
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
	mux.HandleFunc("/v1/recommender/request", sendRecommendRequest)
	mux.HandleFunc("/v1/recommender/response", getRecommendResponse)

	n := negroni.Classic()
	n.UseHandler(mux)
	http.ListenAndServe(":3000", n)
}
