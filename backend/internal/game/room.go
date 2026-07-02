package game

import (
	"sync"
	"time"
	"fmt"
)

type GameState string

const (
	Waiting  GameState = "WAITING"
	Running  GameState = "RUNNING"
	Finished GameState = "FINISHED"
)

type Room struct {

	// Unique Room ID
	ID string

	// Players in this room
	Players map[string]*Player

	// Ball used in the game
	Ball *Ball

	// Current game state
	State GameState

	// Maximum allowed players
	MaxPlayers int

	// Room creation time
	CreatedAt time.Time

	// Protects shared room data
	Mutex sync.RWMutex
}

func NewRoom(id string, maxPlayers int) *Room {

	return &Room{
		ID:         id,
		Players:    make(map[string]*Player),
		Ball:       NewBall(),
		State:      Waiting,
		MaxPlayers: maxPlayers,
		CreatedAt:  time.Now(),
	}
}
func (r *Room) AddPlayer(player *Player) error {

	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	// Game already started
	if r.State != Waiting {
		return fmt.Errorf("game already started")
	}

	// Room full
	if len(r.Players) >= r.MaxPlayers {
		return fmt.Errorf("room is full")
	}

	// Duplicate player
	if _, exists := r.Players[player.ID]; exists {
		return fmt.Errorf("player already exists")
	}

	r.Players[player.ID] = player

	return nil
}

func (r *Room) RemovePlayer(playerID string) {

	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	delete(r.Players, playerID)
}

func (r *Room) GetPlayer(playerID string) (*Player, bool) {

	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	player, exists := r.Players[playerID]

	return player, exists
}

func (r *Room) PlayerCount() int {

	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	return len(r.Players)
}
func (r *Room) CanStart() bool {

	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	return r.State == Waiting &&
		len(r.Players) >= 2
}

func (r *Room) Start() {

	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	r.State = Running
}

func (r *Room) Finish() {

	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	r.State = Finished
}