package websocket 
import (
	"log"
	"pingpong_game/backend/internal/events"
	"pingpong_game/backend/internal/game"
)


type Hub struct{
	Clients map[string]*Client
	Register chan *Client
	Unregister chan *Client 
	Broadcast chan []byte
	room *game.Room
}


func NewHub(room *game.Room) *Hub{
	return &Hub {
		Clients:make(map[string]*Client),
		Register: make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast: make(chan []byte),
		room: room,

	}
}

func (h *Hub) Run(){
	for {

	select {

	case client := <-h.Register:

		h.Clients[client.ID] = client

		log.Printf(
			"Client Connected : %s | Total : %d",
			client.ID,
			len(h.Clients),
		)

	case client := <-h.Unregister:

		if _, ok := h.Clients[client.ID]; ok {

			delete(h.Clients, client.ID)

			close(client.Send)

			log.Printf(
				"Client Disconnected : %s | Total : %d",
				client.ID,
				len(h.Clients),
			)
		}

	case message := <-h.Broadcast:
		log.Println("📤 Hub received broadcast message")

		for _, client := range h.Clients {
			log.Printf("📨 Sending message to client %s", client.ID)


			select {

			case client.Send <- message:

			default:

				close(client.Send)

				delete(h.Clients, client.ID)
			}
		}
	}
}
}
func (h *Hub) Handle(event events.Event) {

	h.BroadcastRoomState(h.room)
}
