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

	log.Printf("waiting for message. to exit press CTRL+C")

	forever := make(chan bool)
	// EXPERIMENT
	chP, err := conn.Channel()
	errorWrapper(err, "Failed to open a channel")
	defer chP.Close()
	err = chP.ExchangeDeclare(
		"livetest-product1",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println(err.Error())
		return
	}
	qP, err := chP.QueueDeclare(
		"livetest-product1", //name
		false,               // durable
		false,               //delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	errorWrapper(err, "Failed to declare a queue")
	err = chP.QueueBind(
		qP.Name,
		"theroute2",
		"livetest-product1",
		false,
		nil,
	)
	if err != nil {
		log.Println(err.Error())
		return
	}
	msgP, err := chP.Consume(
		qP.Name, // queue
		"",      // consumer
		true,    // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	errorWrapper(err, "Failed to register a consumer")

	for d := range msgP {
		go func() {
			log.Printf("received as message: %s\nfrom queue: %s", string(d.Body), qP.Name)
		}()
	}

	<-forever
}
