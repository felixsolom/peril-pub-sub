package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
)

func HandlerMove(gs *gamelogic.GameState) func(gamelogic.ArmyMove) pubsub.Acktype {
	return func(move gamelogic.ArmyMove) pubsub.Acktype {
		defer fmt.Print("> ")
		moveOutcome := gs.HandleMove(move)
		if moveOutcome == gamelogic.MoveOutComeSafe || moveOutcome == gamelogic.MoveOutcomeMakeWar {
			return pubsub.Ack
		}
		return pubsub.NackDiscard
	}
}
