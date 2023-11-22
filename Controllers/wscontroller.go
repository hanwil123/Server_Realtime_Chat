package Controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go-chat/Models"
	"net/http"
)

type Handler struct {
	hub *Hub
}

func NewHandler(h *Hub) *Handler {
	return &Handler{
		hub: h,
	}
}
func (h *Handler) CreateRoom(c *gin.Context) {
	var req Models.CreateRoomReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		panic(err)
	}
	h.hub.Rooms[req.ID] = &Room{
		ID:      req.ID,
		Name:    req.Name,
		Clients: make(map[string]*Client),
	}
	RDB.Create(&req)
	c.JSON(http.StatusOK, req)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) JoinRoom(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roomID := c.Param("roomId")
	username := c.Query("username")
	clientID := c.Query("userId")
	if err != nil {
		// Handle error
	}
	cl := &Client{
		Conn:     conn,
		Message:  make(chan *Message, 10),
		ID:       clientID,
		RoomID:   roomID,
		Username: username,
	}
	m := &Message{
		Content:  "A new user has been joined",
		RoomID:   roomID,
		Username: username,
	}
	h.hub.Register <- cl
	h.hub.Broadcast <- m
	go cl.writeMessage()
	cl.readMessage(h.hub)
}

type RoomRes struct {
	ID     string `gorm:"primaryKey" json:"id"`
	RoomID string `gorm:"foreignKey" json:"roomId"`
	Name   string `json:"name"`
}

func (h *Handler) GetRooms(c *gin.Context) {
	var rooms []Models.CreateRoomReq
	RDB.Find(&rooms)

	roomRes := make([]RoomRes, 0)
	for _, r := range rooms {
		roomRes = append(roomRes, RoomRes{
			ID:     r.ID,
			RoomID: r.RoomID,
			Name:   r.Name,
		})
	}
	c.JSON(http.StatusOK, roomRes)
}

type ClientRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (h *Handler) GetClients(c *gin.Context) {
	var clients []ClientRes
	roomId := c.Param("roomId")
	room, roomExists := h.hub.Rooms[roomId]
	if !roomExists {
		clients = make([]ClientRes, 0)
		c.JSON(http.StatusOK, clients)
		return
	}
	if room.Clients == nil {
		fmt.Printf("Clients pada room tersebut tidak ada")
		c.JSON(http.StatusOK, clients)
		return
	}
	for _, client := range room.Clients {
		clients = append(clients, ClientRes{
			ID:       client.ID,
			Username: client.Username,
		})
	}
	c.JSON(http.StatusOK, clients)
}
