package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"github.com/rabbitmq/amqp091-go"
)

func HandlerMove(gs *gamelogic.GameState, ch *amqp091.Channel) func(gamelogic.ArmyMove) pubsub.Acktype {
	return func(move gamelogic.ArmyMove) pubsub.Acktype {
		defer fmt.Print("> ")

		moveOutcome := gs.HandleMove(move)
		if moveOutcome == gamelogic.MoveOutComeSafe {
			return pubsub.Ack
		}
		if moveOutcome == gamelogic.MoveOutcomeMakeWar {
			userName := gs.GetUsername()
			routingKey := fmt.Sprintf("%s.%s", routing.WarRecognitionsPrefix, userName)
			err := pubsub.PublishJSON(ch,
				routing.ExchangePerilTopic,
				routingKey,
				moveOutcome)
			if err != nil {
				fmt.Printf("Failed to publish war recognition", err)
			}
			return pubsub.NackRequeue
		}
		return pubsub.NackDiscard
	}
}
