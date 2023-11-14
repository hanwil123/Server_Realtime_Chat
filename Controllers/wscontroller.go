package Controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go-chat/Databases"
	"go-chat/Models"
	"net/http"
)

type Handler struct {
	hub *Models.Hub
}

func NewHandler(h *Models.Hub) *Handler {
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
	h.hub.Rooms[req.ID] = &Models.Room{
		ID:      req.ID,
		Name:    req.Name,
		Clients: make(map[string]*Models.Client),
	}
	Databases.RDB.Create(&req)
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
	var req Models.UserClients
	if err := c.ShouldBindJSON(&req); err != nil {
		conn.Close() // Close the WebSocket connection on error
		return
	}
	roomID := c.Param("roomId")
	cl := &Models.Client{
		Socket:   conn,
		Message:  make(chan *Models.Message, 10),
		RoomID:   roomID,
		UserName: req.Username,
		UserID:   req.UserID,
	}
	m := &Models.Message{
		Content:  "A new user has been joined",
		RoomID:   roomID,
		UserName: req.Username,
	}
	h.hub.Register <- cl
	h.hub.Broadcast <- m
	go cl.WriteMessage()
	go cl.ReadMessage(h.hub)
}
func (h *Handler) CreateClients(c *gin.Context) {
	var req Models.UserClients
	errr := c.ShouldBindJSON(&req)
	if errr != nil {
		panic(errr)
	}

	// Create new user
	user := &Models.UserClients{
		Username: req.Username,
		UserID:   req.UserID,
	}
	result := Databases.CDB.Create(user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

type RoomRes struct {
	ID     string `gorm:"primaryKey" json:"id"`
	RoomID string `gorm:"foreignKey" json:"roomId"`
	Name   string `json:"name"`
}

func (h *Handler) GetRooms(c *gin.Context) {
	var rooms []Models.CreateRoomReq
	Databases.RDB.Find(&rooms)

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
	if _, ok := h.hub.Rooms[roomId]; !ok {
		clients = make([]ClientRes, 0)
		c.JSON(http.StatusOK, clients)
	}
	for _, c := range h.hub.Rooms[roomId].Clients {
		clients = append(clients, ClientRes{
			ID:       c.ID,
			Username: c.UserName,
		})
	}
	c.JSON(http.StatusOK, clients)
}
