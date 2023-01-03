package kafka

import (
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaWriteHandler struct {
	BootstrapServer string
}

// create new Kafka write handler
func NewKafkaWriteHandler(bootstrapServer string) *KafkaWriteHandler {
	return &KafkaWriteHandler{BootstrapServer: bootstrapServer}
}

// write message to Kafka queue
func (k *KafkaWriteHandler) WriteMessage(message string) error {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal("no hostname, error = ", err)
		os.Exit(1)
	}
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": k.BootstrapServer,
		"client.id":         hostname,
		"acks":              "all",
	})

	if err != nil {
		log.Fatal("error initializting kafka producer, error = ", err)
		return err
	} else {
		log.Println("kafka producer initialized ", p)
	}

	return nil
}
