package websocket


import (
	"encoding/json"
	"log"
	"time"

	"pingpong_game/backend/internal/events"

	"github.com/gorilla/websocket"
)

func (c *Client) ReadPump() {

	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	// Maximum size of incoming message
	c.Conn.SetReadLimit(1024)

	for {

		// Read message from browser
		_, data, err := c.Conn.ReadMessage()
		log.Printf("Received Raw JSON: %s", string(data))
		if err != nil {

			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Println("Read Error:", err)
			}

			break
		}

		// Convert JSON into Message struct
		var msg Message

		err = json.Unmarshal(data, &msg)
		if err != nil {
			log.Println("Invalid JSON:", err)
			continue
		}

		// Convert websocket message into game event
		event := events.Event{
			Type:           events.EventType(msg.Type),
			PlayerID:       msg.PlayerID,
			PlayerName:     msg.PlayerName,
			TargetPlayerID: msg.TargetPlayerID,
			TimeStamp:      time.Now(),
		}

		// Publish event to game engine
		c.Engine.Publish(event)
		log.Printf("Published Event: %+v", event)
	}
}