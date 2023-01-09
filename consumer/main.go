package main

import (
	"log"

	"github.com/streadway/amqp"
)

func errorWrapper(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	errorWrapper(err, "Failed to connect rabbitmq")
	defer conn.Close()

	ch, err := conn.Channel()
	errorWrapper(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"golang-queue", //name
		false,          // durable
		false,          //delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	errorWrapper(err, "Failed to declare a queue")

	msg, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	errorWrapper(err, "Failed to register a consumer")

	forever := make(chan bool)
	go func() {
		for d := range msg {
			log.Printf("received as message: %s\nfrom queue: %s", string(d.Body), q.Name)
		}
	}()

	qP, err := ch.QueueDeclare(
		"livetest-product1", //name
		false,               // durable
		false,               //delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	errorWrapper(err, "Failed to declare a queue")

	msgP, err := ch.Consume(
		qP.Name, // queue
		"",      // consumer
		true,    // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	errorWrapper(err, "Failed to register a consumer")

	go func() {
		for d := range msgP {
			log.Printf("received as message: %s\nfrom queue: %s", string(d.Body), qP.Name)
		}
	}()

	log.Printf("waiting for message. to exit press CTRL+C")
	<-forever
}
