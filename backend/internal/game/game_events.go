package game

import (
	"math/rand"
	"time"
	"pingpong_game/backend/internal/events"
)

func (e *Engine) handleGameFinished(event events.Event) {

}

func (e *Engine) handleGameStarted() {

	if !e.validateGameStart() {
		return
	}

	firstPlayer := e.chooseFirstPlayer()

	e.initializeBall(firstPlayer)
}

func (e *Engine) validateGameStart() bool {

	if e.room.State == Running {
		return false
	}

	if len(e.room.Players) < 2 {
		return false
	}

	return true
}

func (e *Engine) chooseFirstPlayer() *Player {

	players := make([]*Player, 0, len(e.room.Players))

	for _, player := range e.room.Players {
		players = append(players, player)
	}

	index := rand.Intn(len(players))

	return players[index]
}

func (e *Engine) initializeBall(player *Player) {

	e.room.State = Running

	e.room.Ball.CurrentOwnerID = player.ID
	e.room.Ball.PreviousOwnerID = ""
	e.room.Ball.PassCount = 0
	e.room.Ball.ReceivedAt = time.Now()
}