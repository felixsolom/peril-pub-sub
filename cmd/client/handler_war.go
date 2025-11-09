package main

import (
	"fmt"
	"time"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"github.com/rabbitmq/amqp091-go"
)

func HandlerWar(gs *gamelogic.GameState, ch *amqp091.Channel) func(gamelogic.RecognitionOfWar) pubsub.Acktype {
	return func(rw gamelogic.RecognitionOfWar) pubsub.Acktype {
		defer fmt.Print("> ")
		outcome, winner, loser := gs.HandleWar(rw)

		fmt.Printf("War received from %s\n", rw.Attacker.Username)

		switch outcome {
		case gamelogic.WarOutcomeOpponentWon:
			message := fmt.Sprintf("%s won a war against %s", winner, loser)
			return publishGameLog(gs, message, ch)
		case gamelogic.WarOutcomeYouWon:
			message := fmt.Sprintf("%s won a war against %s", winner, loser)
			return publishGameLog(gs, message, ch)
		case gamelogic.WarOutcomeDraw:
			message := fmt.Sprintf("A war between %s and %s resulted in a draw", winner, loser)
			return publishGameLog(gs, message, ch)
		case gamelogic.WarOutcomeNoUnits:
			return pubsub.NackDiscard
		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.NackRequeue
		default:
			fmt.Println("War outcome not known")
			return pubsub.NackDiscard
		}
	}
}

func publishGameLog(gs *gamelogic.GameState, message string, ch *amqp091.Channel) pubsub.Acktype {
	gameLog := routing.GameLog{
		CurrentTime: time.Now(),
		Message:     message,
		Username:    gs.GetUsername(),
	}

	routingKey := fmt.Sprintf("%s.%s", routing.GameLogSlug, gs.GetUsername())
	err := pubsub.PublishGob(
		ch,
		routing.ExchangePerilTopic,
		routingKey,
		gameLog,
	)
	if err != nil {
		fmt.Printf("coudln't publish game log: %v\n", err)
		return pubsub.NackRequeue
	}
	return pubsub.Ack
}
