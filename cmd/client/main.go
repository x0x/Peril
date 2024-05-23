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
	fmt.Println("Starting Peril client...")

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("Unable to connect to rabbitMQ server")
	}

	defer conn.Close()
	fmt.Println("Connection was Successful")

	userName, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Printf("username not available for the player :%v", err)
	}

	_, _, err = pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect,
		routing.PauseKey+"."+userName,
		routing.PauseKey,
		pubsub.SimpleQueueType{
			Durable:    false,
			AutoDelete: true,
			Exclusive:  true,
			NoWait:     false,
			Args:       nil,
		},
	)
	if err != nil {
		log.Fatalf("Unable to declare and bind the queue :%v", err)
	}

	gs := gamelogic.NewGameState(userName)

	pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+userName,
		routing.PauseKey,
		pubsub.TransientQueueType,
    gamelogic.HandlerPause(gs),
	)

	pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilTopic,
		routing.ArmyMovesPrefix+"."+userName,
		routing.ArmyMovesPrefix+".*",
		pubsub.TransientQueueType,
		gamelogic.HandlerMove(gs),
	)

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	for {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}

		switch input[0] {

		case "spawn":
			err := gs.CommandSpawn(input)
			if err != nil {
				fmt.Printf("Error in CommandSpawn :%v\n", err)
			}

		case "move":
			move, err := gs.CommandMove(input)
			if err != nil {
				fmt.Printf("Error in CommandMove :%+v\n", err)
				continue
			}

			pubsub.PublishJSON(
				ch,
				routing.ExchangePerilTopic,
				routing.ArmyMovesPrefix+"."+userName,
				move,
			)
		case "status":
			gs.CommandStatus()

		case "help":
			gamelogic.PrintClientHelp()

		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			fmt.Println("Command Not Identified")
		}
	}
}
