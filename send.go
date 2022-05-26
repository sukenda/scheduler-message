package main

import (
	"encoding/json"
	"fmt"
	"github.com/sukenda/scheduler-message/config"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ip := os.Getenv("RABBITMQ_IP")
	port := os.Getenv("RABBITMQ_PORT")
	conn, err := amqp.Dial(fmt.Sprintf("amqp://guest:guest@%v:%v/", ip, port))
	FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	c := config.Config{
		Queue:    "delayed-exchange-queue",
		Key:      "delayed-key",
		Exchange: "delayed-exchange",
		Durable:  true,
	}

	// Declare exchange
	args := make(amqp.Table)
	args["x-delayed-type"] = "direct"
	err = ch.ExchangeDeclare("delayed-exchange", "x-delayed-message", true, false, false, false, args)
	if err != nil {
		log.Fatalf(err.Error())
	}

	_, err = ch.QueueDeclare(
		c.Queue,
		c.Durable,
		false,
		false,
		false,
		nil,
	)

	err = ch.QueueBind(c.Queue, c.Key, c.Exchange, false, nil)
	if err != nil {
		log.Fatalf(err.Error())
	}

	payload := config.Payload{
		ID:      uuid.New().String(),
		Message: fmt.Sprintf("Message delay with %v", "Value"),
		Time:    time.Now(),
	}

	bytes, _ := json.Marshal(payload)

	delay := 1000 * 60 * 10
	headers := make(amqp.Table)
	headers["x-delay"] = delay

	err = ch.Publish(
		c.Exchange,
		c.Key,
		false,
		false,
		amqp.Publishing{
			MessageId:    payload.ID,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			ContentType:  "application/bytes",
			Body:         bytes,
			Headers:      headers,
		})
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf(" Sent %s\n", fmt.Sprintf("Publish with message id %v and delay %v", payload.ID, delay))
}
