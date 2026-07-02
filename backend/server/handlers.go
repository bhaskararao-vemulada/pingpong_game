package server

import (
	"fmt"
	"log"
	"net/http"

	"pingpong_game/backend/internal/websocket"

	gorilla "github.com/gorilla/websocket"
)

var upgrader = gorilla.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	// Allow all origins for now.
	// Later we'll restrict this.
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Server is running 🚀")
}

func (s *Server) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	

	log.Println(">>> WebSocket handler called")

	

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket Upgrade Error:", err)
		return
	}

	log.Println("New WebSocket Connection")

	// Temporary client id
	clientID := r.RemoteAddr

	// Create websocket client
	client := websocket.NewClient(
	clientID,
	conn,
	s.hub,
	s.engine,
)
	// Register with Hub
	s.hub.Register <- client

	// Start client goroutines
	go client.ReadPump()
	go client.WritePump()
}