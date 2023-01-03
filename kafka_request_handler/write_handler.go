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
func (k *KafkaWriteHandler) WriteMessage(message string, topic string) error {
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

	delivery_channel := make(chan kafka.Event, 10000)
	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: []byte(message),
	}, delivery_channel)
	if err != nil {
		log.Fatal("error sending message, error = ", err)
	} else {
		log.Println("message sent -> ", message)
	}
	e := <-delivery_channel
	m := e.(*kafka.Message)
	if m.TopicPartition.Error != nil {
		log.Fatal(m.TopicPartition.Error)
	} else {
		log.Println("message delivered")
	}

	return nil
}
