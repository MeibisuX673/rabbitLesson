package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
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

	if err := ch.ExchangeDeclare("logs", "fanout", true, false, false, false, nil); err != nil {
		log.Panic(err)
	}

	q, err := ch.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		log.Panic(err)
	}

	if err = ch.QueueBind(q.Name, "", "logs", false, nil); err != nil {
		log.Panic(err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Panic(err)
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Println(string(d.Body))
		}
	}()

	<-forever

}
