package pubsub

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Acktype int

const (
	Ack Acktype = iota
	NackRequeue
	NackDiscard
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
	handler func(T) Acktype,
) error {
	ch, _, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}
	go func() {
		for msg := range msgs {
			var body T
			err := json.Unmarshal(msg.Body, &body)
			if err != nil {
				fmt.Printf("Failed to unmarshal message: %v", err)
				msg.Nack(false, false)
			}
			returnType := handler(body)
			switch returnType {
			case Ack:
				msg.Ack(false)
				fmt.Println("Message processed")
			case NackRequeue:
				msg.Nack(false, true)
				fmt.Println(("Message didn't process. Was requeued"))
			case NackDiscard:
				msg.Nack(false, false)
				fmt.Println("Message didn't process and was discarded")
			default:
				msg.Nack(false, false)
				fmt.Println("Unknown ack command")
			}
		}
	}()
	return nil
}
