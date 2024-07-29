package main

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	topic := "HVSE"

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  "localhost:9092",
		"group.id":           "Data Team",
		"auto.offset.reset":  "smallest",
		"enable.auto.commit": "true",
		"auto.commit.interval.ms": "1000",
	})	
	if err != nil {
		log.Fatal(err)
	}
	err = consumer.Subscribe(topic, nil)
	if err != nil {
		log.Fatal(err)
	}

	for {
		ev := consumer.Poll(100)
		switch e := ev.(type) {
		case *kafka.Message:
			fmt.Printf("Consumed order by Data Team: %s\n", string(e.Value))
			_, err := consumer.CommitMessage(e)
			if err != nil {
				log.Printf("Failed to commit message: %s\n", err)
			}
		case *kafka.Error:
			fmt.Printf("%s\n", e)
		}
	}
}
