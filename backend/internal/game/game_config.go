package game

import (
	"log"
	"time"

	"pingpong_game/backend/internal/events"
)

func (e *Engine) handleGameConfigured(
	event events.Event,
) {
	// 1. Validate duration
	if event.GameDurationSeconds <= 0 {
		log.Printf(
			"❌ Invalid game duration: %d seconds",
			event.GameDurationSeconds,
		)
		return
	}

	// 2. Lock room before modifying shared state
	e.room.Mutex.Lock()

	// 3. Configuration allowed only
	// while room is WAITING
	if e.room.State != Waiting {
		currentState := e.room.State

		e.room.Mutex.Unlock()

		log.Printf(
			"⚠️ Cannot configure game while state is: %s",
			currentState,
		)
		return
	}

	// 4. Store selected duration
	e.room.GameDuration =
		time.Duration(
			event.GameDurationSeconds,
		) * time.Second

	e.room.Mutex.Unlock()

	log.Printf(
		"⚙️ Game configured | Duration: %d seconds",
		event.GameDurationSeconds,
	)

	// 5. Broadcast updated room state
	e.bus.Publish(events.Event{
		Type:      events.EventRoomUpdated,
		TimeStamp: time.Now(),
	})
}