package game 
import ("pingpong_game/backend/internal/events"
"log")

func (e *Engine) handleBallPassed(event events.Event) {

	log.Printf(
		"🏓 Pass requested: %s -> %s",
		event.PlayerID,
		event.TargetPlayerID,
	)

	// Only the current ball owner can pass.
	if e.room.Ball.CurrentOwnerID != event.PlayerID {
		log.Printf(
			"❌ Pass rejected: player %s does not own the ball",
			event.PlayerID,
		)
		return
	}

	// Target player must exist in the room.
	if _, exists := e.room.Players[event.TargetPlayerID]; !exists {
		log.Printf(
			"❌ Pass rejected: target player %s not found",
			event.TargetPlayerID,
		)
		return
	}

	e.room.Ball.PassTo(event.TargetPlayerID)

	log.Printf(
		"✅ Ball passed: %s -> %s | PassCount: %d",
		event.PlayerID,
		event.TargetPlayerID,
		e.room.Ball.PassCount,
	)

	e.bus.Publish(events.Event{
		Type: events.EventRoomUpdated,
	})
}