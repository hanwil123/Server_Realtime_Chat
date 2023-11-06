package main

import (
	"go-chat/Controllers"
	"go-chat/Databases"
	"go-chat/Models"
	"go-chat/Routes"
)

func main() {
	hub := Models.NewHub()
	Databases.ConnectRoomList()
	Databases.ConnectClientUser()
	wsHandler := Controllers.NewHandler(hub)
	Routes.WsRouter(wsHandler)
	go hub.Run()
	Routes.Start("localhost:8080")
}
