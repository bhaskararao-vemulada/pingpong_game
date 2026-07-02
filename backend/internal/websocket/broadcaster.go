package websocket 


import (
	"encoding/json"
	"log"

	"pingpong_game/backend/internal/game"
)

func (h *Hub) BroadcastRoomState(room *game.Room) {

	state := game.BuildRoomState(room)

	data, err := json.Marshal(state)
	if err != nil {
		log.Println("Failed to marshal room state:", err)
		return
	}

	h.Broadcast <- data
}