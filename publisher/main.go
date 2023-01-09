package main

import (
	"encoding/json"
	"fmt"
	"log"

	"rabbitmq-go/model"

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

	for _, msg := range msgQueue {
		publishMsg(ch, q, &msg)
		msgFormated := model.BodyPublishTest{
			ProductName:    "livetest-product1",
			HandlerNameKey: "LOC_VER_V3_BS_Get_Score_V3",
			MsgQue:         msg,
		}
		q.Name = msgFormated.ProductName
		publishMsgJson(ch, q, &msgFormated)
	}
}

func publishMsg(ch *amqp.Channel, q amqp.Queue, msg *model.LivetestBacktestMessageQueue) {
	err := ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immadiate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg.ID),
		})

	errorWrapper(err, "Failed to publish message")
	log.Printf("Sending message success: %s", msg.ID)
}

// jika kita send json, lebih baik sehingga di dalam handler queue tidak membutuhkan query get
func publishMsgJson(ch *amqp.Channel, q amqp.Queue, msg *model.BodyPublishTest) {
	body, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immadiate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		})

	errorWrapper(err, "Failed to publish message")
	log.Printf("Sending message success: %s", string(body))
}

var data = model.BatchUploadCSV{
	ID: "1234-4321",
}
var msgQueue = []model.LivetestBacktestMessageQueue{
	{
		ID:               "abcd-dcba-1",
		BatchUploadCSVID: "1234-4321",
	},
	{
		ID:               "abcd-dcba-2",
		BatchUploadCSVID: "1234-4321",
	},
}
