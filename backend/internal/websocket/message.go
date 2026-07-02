package websocket



type Message struct {
	Type           string `json:"type"`
	PlayerID       string `json:"playerId"`
	PlayerName     string `json:"playerName"`
	TargetPlayerID string `json:"targetPlayerId"`
}