package events

import (
	"fmt"
	"sync"

)

type Bus struct {
	subscribers map[EventType][]Subscriber
	mutex       sync.RWMutex
}

func NewBus() *Bus {
	return &Bus{
		subscribers: make(map[EventType][]Subscriber),
	}
}

func (b *Bus) Subscribe(eventType EventType, subscriber Subscriber) error {

	b.mutex.Lock()
	defer b.mutex.Unlock()

	subs := b.subscribers[eventType]

	// Duplicate check
	for _, sub := range subs {

		if sub == subscriber {
			return fmt.Errorf(
				"subscriber already registered for event %s",
				eventType,
			)
		}
	}

	b.subscribers[eventType] = append(subs, subscriber)

	return nil
}

func (b *Bus) Publish(event Event) {

	// Step 1: Read subscribers safely
	b.mutex.RLock()

	subscribers := append(
		[]Subscriber(nil),
		b.subscribers[event.Type]...,
	)

	b.mutex.RUnlock()

	// Step 2: Notify every subscriber
	for _, subscriber := range subscribers {

		go func(sub Subscriber) {

			defer func() {
				if r := recover(); r != nil {
					fmt.Printf(
						"subscriber panic while handling %s : %v\n",
						event.Type,
						r,
					)
				}
			}()

			sub.Handle(event)

		}(subscriber)
	}
}

