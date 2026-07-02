package game

type Notifier interface {
	BroadcastRoomState(room *Room)
}