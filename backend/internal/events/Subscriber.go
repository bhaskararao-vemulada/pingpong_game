package events





type Subscriber interface {
	Handle(event Event)
}