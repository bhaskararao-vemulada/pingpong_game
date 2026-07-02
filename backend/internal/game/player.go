package game 

import (
	
	"time")


type Player struct {
	ID string 
	Name string 
	Avatar string 
	Connected bool
	TotalPosessionTime time.Duration
	PassesMade int 
	PassesRecieved int 
	JoinedAt time.Time 
}


func NewPlayer(id,name string)*Player{
	return &Player{
		ID: id,
		Name: name,

		Connected: true,
		JoinedAt: time.Now(),
		PassesMade: 0,
		PassesRecieved: 0,
		TotalPosessionTime: 0,
	}

}