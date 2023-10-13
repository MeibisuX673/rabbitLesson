package main

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"strings"
	"time"
)

func BodyFromForTopic(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}

func main() {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Panic(err)
	}
	defer ch.Close()

	if err := ch.ExchangeDeclare("topic_log", "topic", true, false, false, false, nil); err != nil {
		log.Panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	body := BodyFromForTopic(os.Args)
	if err := ch.PublishWithContext(ctx,
		"topic_log",
		severityFromForTopic(os.Args),
		false,
		false,
		amqp.Publishing{ContentType: "text/plain", Body: []byte(body)},
	); err != nil {
		log.Panic(err)
	}
}

func severityFromForTopic(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "anonymous.info"
	} else {
		s = os.Args[1]
	}
	return s
}
