package main

import (
    "encoding/json"
    "fmt"
    // "log"

    "github.com/IBM/sarama"
    "github.com/gofiber/fiber/v2"
)

type Comment struct {
    Text string `form:"text" json:"text"`
}

func main() {
    app := fiber.New()
    api := app.Group("/api/v1")
    api.Post("/comments", createComment)
    app.Listen(":3000")
}

func createComment(c *fiber.Ctx) error {
    cmt := new(Comment)
    if err := c.BodyParser(cmt); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "Invalid JSON",
        })
    }

    cmtInBytes, err := json.Marshal(cmt)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Failed to marshal comment",
        })
    }

    if err := PushCommentToQueue("comments", cmtInBytes); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Failed to push comment to queue",
        })
    }

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "success": true,
        "message": "Comment created",
        "comment": cmt,
    })
}

func ConnectProducer(brokersUrl []string) (sarama.SyncProducer, error) {
    config := sarama.NewConfig()
    config.Producer.Return.Successes = true
    config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

    conn, err := sarama.NewSyncProducer(brokersUrl, config)
    if err != nil {
        return nil, err
    }
    return conn, nil
}

func PushCommentToQueue(topic string, message []byte) error {
    brokerUrl := []string{"localhost:29092"}
    producer, err := ConnectProducer(brokerUrl)
    if err != nil {
        return err
    }
    defer producer.Close()

    msg := &sarama.ProducerMessage{
        Topic: topic,
        Value: sarama.StringEncoder(message),
    }
    partition, offset, err := producer.SendMessage(msg)
    if err != nil {
        return err
    }
    fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
    return nil
}