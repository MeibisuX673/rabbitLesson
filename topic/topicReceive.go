package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
)

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

	q, err := ch.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		log.Panic(err)
	}

	if len(os.Args) < 2 {
		log.Printf("Usage: %s [binding_key]...", os.Args[0])
		os.Exit(0)
	}

	for _, s := range os.Args[1:] {
		log.Printf("Binding queue %s to exchange %s with routing key %s",
			q.Name, "topic_log", s)
		err = ch.QueueBind(
			q.Name,      // queue name
			s,           // routing key
			"topic_log", // exchange
			false,
			nil)
		if err != nil {
			log.Panic(err)
		}
	}

	var forever chan struct{}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Panic(err)
	}

	go func() {
		for d := range msgs {
			log.Println(string(d.Body))
		}
	}()

	<-forever

}
