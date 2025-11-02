// Package pubsub
package pubsub

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType string

const (
	durable   SimpleQueueType = "durable"
	transient SimpleQueueType = "transient"
)

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("could not open shannel on the connecton %v", err)
	}

	var queue amqp.Queue

	amqpTable := amqp.Table{
		"x-dead-letter-exchange": "peril_dlx",
	}

	if queueType == durable {
		queue, err = ch.QueueDeclare(queueName, true, false, false, false, amqpTable)
		if err != nil {
			return nil, amqp.Queue{}, fmt.Errorf("could not declare queue: %w", err)
		}
	}
	if queueType == transient {
		queue, err = ch.QueueDeclare(queueName, false, true, true, false, amqpTable)
		if err != nil {
			return nil, amqp.Queue{}, fmt.Errorf("could not declare queue: %w", err)
		}
	}
	err = ch.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("could not bind a queue to the exchange: %w", err)
	}
	return ch, queue, nil
}
