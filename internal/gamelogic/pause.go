package gamelogic

import (
	"fmt"

	"github.com/x0x/Peril/internal/pubsub"
	"github.com/x0x/Peril/internal/routing"
)

func (gs *GameState) HandlePause(ps routing.PlayingState) {
	defer fmt.Println("------------------------")
	fmt.Println()
	if ps.IsPaused {
		fmt.Println("==== Pause Detected ====")
		gs.pauseGame()
	} else {
		fmt.Println("==== Resume Detected ====")
		gs.resumeGame()
	}
}

func HandlerPause(gs *GameState) func(routing.PlayingState) pubsub.AckType {
	return func(ps routing.PlayingState) pubsub.AckType {
		defer fmt.Print(">")
		gs.HandlePause(ps)
    return pubsub.Ack
	}
}
