package main

import (
	"fmt"
	"github.com/nsqio/go-nsq"
	"log"
)

func main() {
	// Instantiate a producer.
	config := nsq.NewConfig()
	err := config.Set("auth_secret", "%n&yFA2JD85z^g")
	if err != nil {
		fmt.Printf("Failed to set auth_secret %v", err)
	}
	producer, err := nsq.NewProducer("192.168.1.110:4150", config)
	if err != nil {
		log.Fatal(err)
	}

	messageBody := []byte("hello")
	topicName := "test"

	// Synchronously publish a single message to the specified topic.
	// Messages can also be sent asynchronously and/or in batches.
	err = producer.Publish(topicName, messageBody)
	if err != nil {
		log.Fatal(err)
	}

	// Gracefully stop the producer when appropriate (e.g. before shutting down the service)
	producer.Stop()
}
