package game

import (
	"sync"
	"time"
	"fmt"
	"sort"
	"log"
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
	// When the actual game started
    StartedAt time.Time
	GameDuration time.Duration

// Maximum duration of one game
    
	WinnerID   string
	WinnerName string
}


func NewRoom(
	id string,
	maxPlayers int,
	gameDuration time.Duration,
) *Room {

	if gameDuration <= 0 {
		gameDuration = 2 * time.Minute
	}

	return &Room{
		ID:           id,
		Players:      make(map[string]*Player),
		Ball:         NewBall(),
		State:        Waiting,
		MaxPlayers:   maxPlayers,
		CreatedAt:    time.Now(),
		GameDuration: gameDuration,
		WinnerID:     "",
		WinnerName:   "",
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
	r.StartedAt = time.Now()
}

func (r *Room) Finish() {

	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	r.State = Finished
}
func (r *Room) AssignInitialBallOwner() error {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if len(r.Players) == 0 {
		return fmt.Errorf("cannot assign ball: no players in room")
	}

	for playerID := range r.Players {
		r.Ball.CurrentOwnerID = playerID
		r.Ball.PreviousOwnerID = ""
		r.Ball.ReceivedAt = time.Now()
		r.Ball.PassCount = 0

		return nil
	}

	return fmt.Errorf("cannot assign ball owner")
}
func (r *Room) PassBall(
	fromPlayerID string,
	targetPlayerID string,
) error {

	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	// 1. Validate current owner.
	if r.Ball.CurrentOwnerID != fromPlayerID {
		return fmt.Errorf(
			"player %s does not own the ball",
			fromPlayerID,
		)
	}

	// 2. Validate target player.
	targetPlayer, exists := r.Players[targetPlayerID]
	if !exists {
		return fmt.Errorf(
			"target player %s not found",
			targetPlayerID,
		)
	}

	// 3. Get current player.
	currentPlayer, exists := r.Players[fromPlayerID]
	if !exists {
		return fmt.Errorf(
			"current player %s not found",
			fromPlayerID,
		)
	}

	// 4. Calculate how long current player held the ball.
	heldDuration := time.Since(r.Ball.ReceivedAt)

	// 5. Add possession time to current owner.
	currentPlayer.TotalPosessionTime += heldDuration

	// 6. Update pass statistics.
	currentPlayer.PassesMade++
	targetPlayer.PassesRecieved++

	// 7. Move ball to target.
	r.Ball.PassTo(targetPlayerID)

	return nil
}

func (r *Room) NextPlayerID(currentPlayerID string) (string, error) {

	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	if len(r.Players) < 2 {
		return "", fmt.Errorf("not enough players to select next player")
	}

	playerIDs := make([]string, 0, len(r.Players))

	for playerID := range r.Players {
		playerIDs = append(playerIDs, playerID)
	}

	sort.Strings(playerIDs)

	for index, playerID := range playerIDs {

		if playerID == currentPlayerID {

			nextIndex := (index + 1) % len(playerIDs)

			return playerIDs[nextIndex], nil
		}
	}

	return "", fmt.Errorf(
		"current player %s not found in room",
		currentPlayerID,
	)
}
func (r *Room) FinishGame() error {

	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if r.State != Running {
		return fmt.Errorf("game is not running")
	}

	// Capture the player holding the ball
	// at the exact moment the game finishes.
	if r.Ball != nil && r.Ball.CurrentOwnerID != "" {

		currentOwnerID := r.Ball.CurrentOwnerID

		currentPlayer, exists := r.Players[currentOwnerID]
		if !exists {
			return fmt.Errorf(
				"current ball owner %s not found",
				currentOwnerID,
			)
		}

		// Winner = player holding the ball at finish time.
		r.WinnerID = currentPlayer.ID
		r.WinnerName = currentPlayer.Name

		// Add the final possession interval because
		// this player did not pass before game ended.
		finalPossessionDuration := time.Since(
			r.Ball.ReceivedAt,
		)

		currentPlayer.TotalPosessionTime += finalPossessionDuration

		log.Printf(
			"⏱ Final possession added | Player: %s | Duration: %v",
			currentPlayer.Name,
			finalPossessionDuration,
		)

		log.Printf(
			"🏆 Winner selected | ID: %s | Name: %s",
			r.WinnerID,
			r.WinnerName,
		)
	}

	// Stop the game only after final state is captured.
	r.State = Finished

	return nil
}