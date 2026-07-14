package game
import ("pingpong_game/backend/internal/events"
"log"
"time")


func (e *Engine) handleTimer() {

	e.room.Mutex.RLock()

	// Game is not running, so timer has nothing to do.
	if e.room.State != Running {
		e.room.Mutex.RUnlock()
		return
	}

	// Copy game timing values safely.
	startedAt := e.room.StartedAt
	gameDuration := e.room.GameDuration

	// Copy ball timing values safely.
	if e.room.Ball == nil {
		e.room.Mutex.RUnlock()
		return
	}

	if e.room.Ball.CurrentOwnerID == "" {
		e.room.Mutex.RUnlock()
		return
	}

	currentOwnerID := e.room.Ball.CurrentOwnerID
	receivedAt := e.room.Ball.ReceivedAt
	maxHoldTime := e.room.Ball.MaxHoldTime

	e.room.Mutex.RUnlock()

	// --------------------------------
	// 1. CHECK TOTAL GAME DURATION
	// --------------------------------

	gameElapsed := time.Since(startedAt)

	if gameElapsed >= gameDuration {

		log.Printf(
			"🏁 Game duration completed | Total time: %v",
			gameElapsed,
		)

		e.Publish(events.Event{
			Type:      events.EventGameFinished,
			TimeStamp: time.Now(),
		})

		return
	}

	// --------------------------------
	// 2. CHECK BALL HOLD DURATION
	// --------------------------------

	ballElapsed := time.Since(receivedAt)

	if ballElapsed < maxHoldTime {
		return
	}

	log.Printf(
		"⏰ Ball hold time expired | Owner: %s | Held for: %v",
		currentOwnerID,
		ballElapsed,
	)

	e.Publish(events.Event{
		Type:      events.EventTimerExpired,
		PlayerID:  currentOwnerID,
		TimeStamp: time.Now(),
	})
}

func (e *Engine) handleTimerExpired(event events.Event) {

	log.Printf(
		"⏰ Handling timer expired for player: %s",
		event.PlayerID,
	)

	



	// Check whether this timeout event is still valid.
	e.room.Mutex.RLock()
	currentOwnerID := e.room.Ball.CurrentOwnerID
	e.room.Mutex.RUnlock()

	if currentOwnerID != event.PlayerID {
		log.Printf(
			"⚠️ Ignoring stale timer event for player: %s",
			event.PlayerID,
		)
		return
	}

	// Ask Room for the next player in deterministic order.
	targetPlayerID, err := e.room.NextPlayerID(event.PlayerID)

	if err != nil {
		log.Println("❌ Failed to select next player:", err)
		return
	}

	// Perform the validated, concurrency-safe pass.
	err = e.room.PassBall(
		event.PlayerID,
		targetPlayerID,
	)

	if err != nil {
		log.Println("❌ Timeout pass failed:", err)
		return
	}

		e.bus.Publish(events.Event{
			Type: events.EventRoomUpdated,
		})}