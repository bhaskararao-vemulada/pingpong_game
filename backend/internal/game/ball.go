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
func (b *Ball) PassTo(targetPlayerID string) {
	b.PreviousOwnerID = b.CurrentOwnerID
	b.CurrentOwnerID = targetPlayerID
	b.ReceivedAt = time.Now()
	b.PassCount++
}