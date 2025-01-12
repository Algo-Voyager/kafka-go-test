package main

import (
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/IBM/sarama"
)

func ConnectConsumer(brokers string) (sarama.Consumer, error) {
    config := sarama.NewConfig()
    config.Consumer.Return.Errors = true
    consumer, err := sarama.NewConsumer([]string{brokers}, config)
    if err != nil {
        return nil, err
    }
    return consumer, nil
}

func main() {
    topic := "comments"
    brokerUrl := "localhost:29092"
    consumer, err := ConnectConsumer(brokerUrl)
    if err != nil {
        log.Fatalf("Failed to connect consumer: %s", err)
    }
    partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
    if err != nil {
        log.Fatalf("Failed to connect partition: %s", err)
    }
    defer partitionConsumer.Close()

    fmt.Println("Consumer Connected to Kafka")

    sigchan := make(chan os.Signal, 1)
    signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
    msgCount := 0
    doneCh := make(chan struct{})

	// Go routine to consume messages from the partition
    go func() {
        for {
            select {
            case err := <-partitionConsumer.Errors():
                fmt.Println(err)
            case msg := <-partitionConsumer.Messages():
                msgCount++
                fmt.Printf("Received message count: %d | Topic (%s) | Message(%s)\n", msgCount, msg.Topic, string(msg.Value))
            case <-sigchan:
                fmt.Println("Interrupt is detected")
                doneCh <- struct{}{}
                return
            }
        }
    }()
	// Waits for the goroutine to signal completion
    <-doneCh
    fmt.Println("Processed", msgCount, "messages", "from topic", topic)
}
