package game

import "time"

type Ball struct {

	// Player currently holding the ball
	CurrentOwnerID string

	// Previous player who had the ball
	PreviousOwnerID string

	// When the current player received the ball
	ReceivedAt time.Time

	// Maximum duration a player may hold the ball
	MaxHoldTime time.Duration

	// Total number of passes made
	PassCount int
}

func NewBall() *Ball {
	return &Ball{
		CurrentOwnerID:  "",
		PreviousOwnerID: "",
		ReceivedAt:      time.Now(),
		MaxHoldTime:     10 * time.Second,
		PassCount:       0,
	}
}