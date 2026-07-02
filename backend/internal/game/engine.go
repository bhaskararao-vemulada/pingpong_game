package game

import (
	"context"
	"sync"
	"time"
	"pingpong_game/backend/internal/events"
)

type Engine struct {

	// Current room
	room *Room

	// Event queue
	eventChannel chan events.Event
	bus  *events.Bus
	

	// Engine lifecycle
	ctx    context.Context
	cancel context.CancelFunc

	// Fires every second
	ticker *time.Ticker

	// We'll use this later when shutting down gracefully
	wg sync.WaitGroup
}

func NewEngine(room *Room,bus  *events.Bus) *Engine {

	ctx, cancel := context.WithCancel(context.Background())

	return &Engine{
		room:         room,
		eventChannel: make(chan events.Event, 100),
		ctx:          ctx,
		bus:          bus,
		
		cancel:       cancel,
		ticker:       time.NewTicker(time.Second),
	}
}

func (e *Engine) Run() {

	// Remove this for now.
	// We never call wg.Add(1), so Done() will panic.
	// defer e.wg.Done()

	defer e.ticker.Stop()

	for {

		select {

		case event, ok := <-e.eventChannel:

			if !ok {
				return
			}

			e.handleEvent(event)

		case <-e.ticker.C:

			e.handleTimer()

		case <-e.ctx.Done():

			return
		}
	}
}

// Publish sends an event to the engine.
func (e *Engine) Publish(event events.Event) {
	e.eventChannel <- event
}

func (e *Engine) handleEvent(event events.Event) {

	switch event.Type {

	case events.EventPlayerJoined:
		e.handlePlayerJoined(event)

	case events.EventPlayerLeft:
		e.handlePlayerLeft(event)

	case events.EventGameStarted:
		e.handleGameStarted()

	case events.EventGameFinished:
		e.handleGameFinished(event)

	case events.EventTimerExpired:
		e.handleTimerExpired(event)

	case events.EventBallPassed:
		e.handleBallPassed(event)

	default:
		// Unknown event
	}
}