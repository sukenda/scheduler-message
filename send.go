package main

import (
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	queue, err := ch.QueueDeclare(
		"message.delay", // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	FailOnError(err, "Failed to declare a queue")

	headers := make(amqp.Table)
	headers["x-delay"] = 5000

	body := fmt.Sprintf("Message Kenda 5 Detik")
	err = ch.Publish(
		"message.delay", // exchange
		queue.Name,      //  routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			Timestamp:   time.Now(),
			ContentType: "application/json",
			Body:        []byte(body),
			Headers:     headers,
		})
	FailOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s\n", "Done")
}
