package server

import (
	"log"
	"net/http"
	"pingpong_game/backend/internal/events"
	"pingpong_game/backend/internal/game"
	"pingpong_game/backend/internal/websocket"
	"time"
	
	
)

type Server struct {
	httpServer *http.Server
	mux        *http.ServeMux

	// Core application components
	room   *game.Room
	engine *game.Engine
	hub    *websocket.Hub
	bus    *events.Bus
}
func New() *Server {

	mux := http.NewServeMux()

	room := game.NewRoom(
	"room-1",
	4,
	2*time.Minute,
)

	bus := events.NewBus()

	engine := game.NewEngine(room, bus)

	hub := websocket.NewHub(room)
	
	if err := bus.Subscribe(events.EventRoomUpdated, hub); err != nil {
	log.Fatal(err)
    }

	s := &Server{
		mux:    mux,
		room:   room,
		engine: engine,
		hub:    hub,
		bus:    bus,
		
		httpServer: &http.Server{
			Addr:    ":8080",
			Handler: mux,
		},
	}

	s.registerRoutes()

	go s.engine.Run()
	go s.hub.Run()

	return s
}

func (s *Server) Start() error {
	log.Println("HTTP Server listening on :8080")
	return s.httpServer.ListenAndServe()
}
