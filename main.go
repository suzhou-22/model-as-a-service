package main

import (
	"context"
	"encoding/json"
	"log"
	"model-as-a-service/data"
	"model-as-a-service/property"
	"time"

	"github.com/gookit/config/v2"
	yaml "github.com/gookit/config/v2/yaml"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

type Listener struct {
	Conn     *amqp.Connection
	Ch       *amqp.Channel
	Messages <-chan amqp.Delivery
}

func NewListenser(url string) *Listener {
	conn, err := amqp.Dial(url)
	failOnError(err, "Connection Fail")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	q, err := ch.QueueDeclare(
		"rpc_queue", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	return &Listener{conn, ch, msgs}
}

func (l *Listener) Close() {
	l.Ch.Close()
	l.Conn.Close()
}

func (l *Listener) Listen() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for msg := range l.Messages {
		var message data.Message
		err := json.Unmarshal(msg.Body, &message)
		failOnError(err, "Message Format Error")

		resp := data.Message{Question: message.Question, Answer: "Answer"}
		respBytes, _ := json.Marshal(resp)

		err = l.Ch.PublishWithContext(ctx,
			"",          // exchange
			msg.ReplyTo, // routing key
			false,       // mandatory
			false,       // immediate
			amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: msg.CorrelationId,
				Body:          respBytes,
			})

		failOnError(err, "Fail to send back message")
	}
}

func main() {
	config.WithOptions(config.ParseEnv)
	config.AddDriver(yaml.Driver)
	err := config.LoadFiles("config.yaml")

	var configProperty property.ConfigProperty
	config.Decode(&configProperty)

	if err != nil {
		panic(err)
	}

	listener := NewListenser(configProperty.Amqp.Url)
	defer listener.Close()

	go listener.Listen()

	log.Println(config.String("amqp.url") + " Listening...")

	var forever chan any
	<-forever
}
