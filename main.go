package main

import (
	"apipublisher/publisher"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/joho/godotenv"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}

	kafkaHost := os.Getenv("kafka_host")
	kafkaTopics := []string{os.Getenv("kafka_topic")}
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        kafkaHost,
		"group.id":                 "gotest",
		"session.timeout.ms":       60000,
		"auto.offset.reset":        "earliest",
		"enable.auto.offset.store": false,
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created Consumer %v\n", c)

	err = c.SubscribeTopics(kafkaTopics, nil)
	run := true

	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev := c.Poll(1000)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				if e.Headers != nil {
					fmt.Printf("%% Headers: %v\n", e.Headers)
				}

				err := publisher.Run(string(e.Value))
				if err == nil {
					_, err = c.StoreMessage(e)
					if err != nil {
						fmt.Fprintf(os.Stderr, "%% Error storing offset after message %s:\n",
							e.TopicPartition)
					}
				}
			case kafka.Error:
				// terminate the application if all brokers are down.
				fmt.Fprintf(os.Stderr, "%% Error: %v: %v\n", e.Code(), e)
				if e.Code() == kafka.ErrAllBrokersDown {
					run = false
				}
			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}

	fmt.Printf("Closing consumer\n")
	c.Close()
}
