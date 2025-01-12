# Kafka-Go

This project demonstrates a simple Kafka producer and consumer using Go. The producer sends comments to a Kafka topic, and the consumer reads messages from the topic.

## Prerequisites

- Docker
- Docker Compose
- Go 1.23.4 or later

## Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/Algo-Voyager/kafka-go.git
    cd kafka-go
    ```

2. Install the dependencies:

    ```sh
    go mod tidy
    ```

3. Start Kafka and Zookeeper using Docker Compose:

    ```sh
    docker-compose up -d
    ```

## Usage

### Producer

The producer sends comments to the Kafka topic `comments`.

1. Run the producer:

    ```sh
    go run producer/producer.go
    ```

2. Send a POST request to create a comment:

    ```sh
    curl -X POST http://localhost:3000/api/v1/comments -H "Content-Type: application/json" -d '{"text": "Hello, Kafka!"}'
    ```

### Consumer

The consumer reads messages from the Kafka topic `comments`.

1. Run the consumer:

    ```sh
    go run consumer/consumer.go
    ```

## Project Structure

- [producer.go](http://_vscodecontentref_/0): Contains the code for the Kafka producer.
- [consumer.go](http://_vscodecontentref_/1): Contains the code for the Kafka consumer.
- [go.mod](http://_vscodecontentref_/2): Go module file with dependencies.

## Dependencies

- [github.com/IBM/sarama](https://github.com/IBM/sarama): Kafka client library for Go.
- [github.com/gofiber/fiber/v2](https://github.com/gofiber/fiber): Express-inspired web framework for Go.
