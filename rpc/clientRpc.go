package main

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func rpcFibonachi(n int) (int, error) {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return 0, err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		return 0, err
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return 0, err
	}

	corId := randomString(42)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := ch.PublishWithContext(ctx,
		"",
		"rpc_queue",
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corId,
			ReplyTo:       q.Name,
			Body:          []byte(strconv.Itoa(n)),
		},
	); err != nil {
		return 0, err
	}

	for d := range msgs {
		if corId == d.CorrelationId {
			res, err := strconv.Atoi(string(d.Body))
			if err != nil {
				return 0, err
			}
			return res, nil
		}
	}

	return 0, err

}

func main() {

	n := bodyFromForRpc(os.Args)

	log.Printf(" [x] Requesting fib(%d)", n)
	res, err := rpcFibonachi(n)
	if err != nil {
		log.Panic(err)
	}

	log.Printf(" [.] Got %d", res)

}

func bodyFromForRpc(args []string) int {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "30"
	} else {
		s = strings.Join(args[1:], " ")
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		log.Panic(err)
	}
	return n
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
