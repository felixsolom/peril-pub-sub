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

	fmt.Println("Starting Peril server...")

	userName, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Could not produce client welcome sequence: %v", err)
	}

	queueName := routing.PauseKey + "." + userName
	ch, queue, err := pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, queueName, routing.PauseKey, "transient")
	if err != nil {
		log.Fatalf("Could not declare and bind queue to exchance: %v", err)
	}
	defer ch.Close()

	log.Printf("Biding queue: %v, on channel: %v", queue, ch)

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	sig := <-signalChan
	fmt.Println("Received signal:", sig)
	fmt.Println("Shutting down gracefully...")
}
