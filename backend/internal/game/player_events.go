package game
import (
	"log"
	"pingpong_game/backend/internal/events")
func (e *Engine) handlePlayerLeft(event events.Event) {

}
func (e *Engine) handlePlayerJoined(event events.Event) {

	log.Printf("Player Joined: %s", event.PlayerName)

	player := NewPlayer(
		event.PlayerID,
		event.PlayerName,
	)

	if err := e.room.AddPlayer(player); err != nil {
		
		log.Println("Failed to add player:", err)
		return
	}
	e.bus.Publish(events.Event{
	Type: events.EventRoomUpdated,
    })

	log.Printf(
		"Player %s joined successfully. Total Players: %d",
		player.Name,
		e.room.PlayerCount(),
	)
	

	log.Printf("Room State: %+v", BuildRoomState(e.room))


	if e.room.CanStart() {
		log.Println("Enough players joined. Starting game...")

		e.room.Start()

		// We'll publish GAME_STARTED in the next step.
	}
}