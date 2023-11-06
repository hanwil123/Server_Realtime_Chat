package Routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go-chat/Controllers"
)

var r *gin.Engine

func WsRouter(wsHandler *Controllers.Handler) {
	gin.SetMode(gin.ReleaseMode)
	r = gin.Default()
	r.Use(cors.Default())
	r.POST("/ws/createRoom", wsHandler.CreateRoom)
	r.GET("/ws/createRoom", wsHandler.CreateRoom)
	r.GET("/ws/createClients", wsHandler.RegisterClients)
	r.POST("/ws/createClients", wsHandler.RegisterClients)
	r.GET("/ws/joinClients", wsHandler.LoginClients)
	r.POST("/ws/joinClients", wsHandler.LoginClients)
	r.GET("/ws/joinRoom/:roomId", wsHandler.JoinRoom)
	r.POST("/ws/joinRoom/:roomId", wsHandler.JoinRoom)
	r.GET("/ws/GetRoom", wsHandler.GetRooms)
	r.POST("/ws/GetRoom", wsHandler.GetRooms)
	r.GET("/ws/GetClients/:roomId", wsHandler.GetClients)
}

func Start(addr string) error {
	return r.Run(addr)
}
