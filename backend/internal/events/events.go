package events

import "time"

type EventType string

const (
	EventPlayerJoined  EventType = "PLAYER_JOINED"
	EventPlayerLeft    EventType = "PLAYER_LEFT"
	EventBallPassed    EventType = "BALL_PASSED"
	EventTimerExpired  EventType = "TIMER_EXPIRED"
	EventGameStarted   EventType = "GAME_STARTED"
	EventGameFinished  EventType = "GAME_FINISHED"
	EventRoomUpdated   EventType = "ROOM_UPDATED"
	
	EventGameConfigured EventType = "GAME_CONFIGURED"
)

type Event struct {
	Type           EventType
	PlayerID       string
	PlayerName     string
	TargetPlayerID string
	TimeStamp      time.Time
	GameDurationSeconds int `json:"gameDurationSeconds,omitempty"`
}