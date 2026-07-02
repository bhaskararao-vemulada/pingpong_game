package websocket


import (
	"log"

	"github.com/gorilla/websocket"
)

func (c *Client) WritePump() {

	defer func() {
		c.Conn.Close()
	}()

	for {

		select {

		case message, ok := <-c.Send:

			if !ok {

				// Hub closed the channel
				c.Conn.WriteMessage(
					websocket.CloseMessage,
					[]byte{},
				)

				return
			}
            log.Printf("📡 Writing to websocket: %s", string(message))
			err := c.Conn.WriteMessage(
				websocket.TextMessage,
				message,
			)

			if err != nil {
				log.Println("Write Error:", err)
				return
			}
		}
	}
}