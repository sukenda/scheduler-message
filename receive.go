package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sukenda/scheduler-message/config"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func FailOnErr(err error, msg string) {
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
	FailOnErr(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	FailOnErr(err, "Failed to open a channel")
	defer ch.Close()

	c := config.Config{
		Queue:    "delayed-exchange-queue",
		Key:      "delayed-key",
		Exchange: "delayed-exchange",
		Durable:  true,
	}

	_, err = ch.QueueDeclare(
		c.Queue,
		c.Durable,
		false,
		false,
		false,
		nil,
	)
	FailOnErr(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		c.Queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	FailOnErr(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
