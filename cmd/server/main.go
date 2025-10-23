package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	connString := "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(connString)
	if err != nil {
		log.Printf("Couldn't make a new connection: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Could not open connection channel: %v", err)
	}
	defer ch.Close()

	fmt.Println("AMQP connetion established")
	fmt.Println("Starting Peril server...")

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	sig := <-signalChan
	fmt.Println("Received signal:", sig)
	fmt.Println("Shutting down gracefully...")
}
