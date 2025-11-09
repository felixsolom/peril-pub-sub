package pubsub

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishGob[T any](ch *amqp.Channel, exchange, key string, val T) error {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(val)
	if err != nil {
		return err
	}

	err = ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{
		ContentType: "application/gob",
		Body:        buff.Bytes(),
	})
	if err != nil {
		return fmt.Errorf("basic publish: %w", err)
	}
	return nil
}
