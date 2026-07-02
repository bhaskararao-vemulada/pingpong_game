package main

import (
	"log"
	"pingpong_game/backend/server")

func main(){
	log.Println("🏓 Starting Ping Pong Server...")
	srv:=server.New()
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}

}