package websocket

import (
	"pingpong_game/backend/internal/game"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Send chan []byte

	Hub    *Hub
	Engine *game.Engine
}


func NewClient(
	id string,
	conn *websocket.Conn,
	hub *Hub,
	engine *game.Engine,
) *Client {

	return &Client{
		ID:     id,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		Hub:    hub,
		Engine: engine,
	}
}