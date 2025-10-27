package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	connString := "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(connString)
	if err != nil {
		log.Fatalf("Couldn't make a new connection: %v", err)
	}
	defer conn.Close()

	fmt.Println("AMQP connection established")
	fmt.Println("Starting Peril server...")
	gamelogic.PrintServerHelp()

	ch, _, err := pubsub.DeclareAndBind(conn, routing.ExchangePerilTopic, routing.GameLogSlug, "game_logs.*", "durable")
	if err != nil {
		log.Fatalf("could not declare and bind queue to exchange: %v", err)
	}
	defer ch.Close()

inputLoop:
	for {
		line := gamelogic.GetInput()
		switch line[0] {
		case "pause":
			fmt.Println("sending pause message...")
			err = pubsub.PublishJSON(ch, string(routing.ExchangePerilDirect), string(routing.PauseKey), routing.PlayingState{
				IsPaused: true,
			})
			if err != nil {
				log.Fatalf("Could not publish json: %v", err)
			}
		case "resume":
			fmt.Println("sending resume message...")
			err = pubsub.PublishJSON(ch, string(routing.ExchangePerilDirect), string(routing.PauseKey), routing.PlayingState{
				IsPaused: false,
			})
			if err != nil {
				log.Fatalf("Could not publish json: %v", err)
			}
		case "quit":
			fmt.Println("Exiting the game")
			break inputLoop
		default:
			fmt.Println("Unknown command")
		}
	}

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	sig := <-signalChan
	fmt.Println("Received signal:", sig)
	fmt.Println("Shutting down gracefully...")
}
