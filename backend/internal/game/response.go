package game 
import( "time"
"math")


type Activity struct {
	Time string `json:"time"`
	Message string `json:"string"`
}

type Response struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type PlayerState struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Connected bool `json:"connected"`
	TotalPossessionTime int64 `json:"totalPossessionTime"`

	PassesMade int `json:"passesMade"`
	PassesRecieved int `json:"passesRecieved"`
}

type RoomState struct {
	RoomID string `json:"roomId"`
	State GameState `json:"state"`

	CurrentOwnerID string `json:"currentOwnerId"`

	CurrentOwnnerName string `json:"currentOwnerName"`

	RemainingSeconds int  `json:"remainingSeconds"`

	PassCount int `json:"passCount"`
	
	Players []PlayerState `json:"playerState"`

	Activities []Activity `json:"activities"`
	

}

type PlayerResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Connected bool   `json:"connected"`

	TotalPossessionTime int64 `json:"totalPossessionTime"`
	PassesMade          int   `json:"passesMade"`
	PassesRecieved      int   `json:"passesRecieved"`
}

type BallResponse struct {
	CurrentOwnerID string `json:"currentOwnerId"`
	PreviousOwnerID string `json:"previousOwnerId"`
	PassCount      int    `json:"passCount"`
}

type RoomStateResponse struct {
	Type       string            `json:"type"`
	RoomID     string            `json:"roomId"`
	GameState  GameState         `json:"gameState"`
	MaxPlayers int               `json:"maxPlayers"`
	Players    []PlayerResponse  `json:"players"`
	Ball       BallResponse      `json:"ball"`
	WinnerID   string `json:"winnerId"`
	WinnerName string `json:"winnerName"`
	RemainingSeconds int `json:"remainingSeconds"`
	GameStartedAt  time.Time `json:"gameStartedAt"`
	GameEndsAt     time.Time `json:"gameEndsAt"`

	BallReceivedAt time.Time `json:"ballReceivedAt"`
	BallExpiresAt  time.Time `json:"ballExpiresAt"`


}

	

func BuildRoomState(room *Room) RoomStateResponse {

	response := RoomStateResponse{
		Type:       "ROOM_STATE",
		RoomID:     room.ID,
		GameState:  room.State,
		MaxPlayers: room.MaxPlayers,
		Players:    make([]PlayerResponse, 0, len(room.Players)),
		WinnerID:   room.WinnerID,
		WinnerName: room.WinnerName,
		GameStartedAt: room.StartedAt,
		GameEndsAt:    room.StartedAt.Add(room.GameDuration),
}
		

	
	

	// Build Players
	for _,player := range room.Players {

	response.Players = append(
		response.Players,
		PlayerResponse{
			ID:        player.ID,
			Name:      player.Name,
			Connected: player.Connected,

			TotalPossessionTime: int64(
				player.TotalPosessionTime / time.Second,
			),

			PassesMade:     player.PassesMade,
			PassesRecieved: player.PassesRecieved,
		},
	)
}

	// Build Ball
	if room.Ball != nil {

		response.Ball = BallResponse{
			CurrentOwnerID:  room.Ball.CurrentOwnerID,
			PreviousOwnerID: room.Ball.PreviousOwnerID,
			PassCount:       room.Ball.PassCount,
		}

	// NEW
		response.BallReceivedAt = room.Ball.ReceivedAt

		response.BallExpiresAt =
			room.Ball.ReceivedAt.Add(
				room.Ball.MaxHoldTime,
			)
	}
	if room.State == Running &&
	room.Ball != nil &&
	room.Ball.CurrentOwnerID != "" {

	elapsed := time.Since(room.Ball.ReceivedAt)

	remaining := room.Ball.MaxHoldTime - elapsed
	if remaining > 0 {
		response.RemainingSeconds = int(
			math.Ceil(remaining.Seconds()),
		)
	} else {
		response.RemainingSeconds = 0
	}
	
}

	return response
}