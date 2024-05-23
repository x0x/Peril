package main

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/x0x/Peril/internal/gamelogic"
	"github.com/x0x/Peril/internal/pubsub"
	"github.com/x0x/Peril/internal/routing"
)

func main() {
	fmt.Println("Starting Peril server...")

	connStr := "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(connStr)
	if err != nil {
		log.Fatal("Unable to connect to rabitMQ server")
	}
	defer conn.Close()
	fmt.Println("Connection was Successful")


  ch, _, err := pubsub.DeclareAndBind(conn, routing.ExchangePerilTopic,
		routing.GameLogSlug,
		"game_logs.*",
		pubsub.SimpleQueueType{
			Durable:    true,
			AutoDelete: false,
			Exclusive:  false,
			NoWait:     false,
			Args:       nil,
		},
	)
	if err != nil {
		log.Fatalf("Unable to declare and bind the queue :%v", err)
	}
	gamelogic.PrintServerHelp()

	for {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}

		switch input[0] {
		case "pause":
			log.Printf("Publishing pause message")
			pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
				IsPaused: true,
			})

		case "resume":
			log.Printf("Publishing resume message")
			pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
				IsPaused: false,
			})

		case "quit":
			log.Printf("Exiting the server ...")
			return

		default:
			log.Print("Could not understand the command !")
		}
	}
}
