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

	if err := ch.ExchangeDeclare("direct_log", "direct", true, false, false, false, nil); err != nil {
		log.Panic(err)
	}

	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		log.Panic(err)
	}

	for _, s := range os.Args[1:] {
		log.Printf("Binding queue %s to exchange %s with routing key %s",
			q.Name, "logs_direct", s)
		err = ch.QueueBind(
			q.Name,       // queue name
			s,            // routing key
			"direct_log", // exchange
			false,
			nil)
		if err != nil {
			log.Panic(err)
		}
	}

	var forever chan struct{}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
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
