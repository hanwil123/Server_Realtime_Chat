package main

import (
	"go-chat/Controllers"
	"go-chat/Routes"
)

func main() {
	hub := Controllers.NewHub()
	Controllers.ConnectRoomList()
	Controllers.ConnectUsers()
	wsHandler := Controllers.NewHandler(hub)
	Routes.WsRouter(wsHandler)
	go hub.Run()
	Routes.Start("localhost:8080")
}
