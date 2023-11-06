package Models

import (
	"github.com/gorilla/websocket"
	"log"
)

type Client struct {
	Socket   *websocket.Conn
	Message  chan *Message
	ID       string `gorm:"AutoIncrement" json:"id"`
	UserID   uint64 `gorm:"foreignKey" json:"userId"`
	RoomID   string `gorm:"foreignKey" json:"roomId"`
	UserName string `json:"username"`
}

type Message struct {
	Content  string `json:"content"`
	RoomID   string `json:"roomId"`
	UserName string `json:"username"`
}

func (c *Client) WriteMessage() {
	defer func() {
		c.Socket.Close()
	}()
	for {
		message, ok := <-c.Message
		if !ok {
			return
		}
		c.Socket.WriteJSON(message)
	}
}

func (c *Client) ReadMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Socket.Close()
	}()

	for {
		_, m, err := c.Socket.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error %v", err)
			}
			break
		}
		msg := &Message{
			Content:  string(m),
			RoomID:   c.RoomID,
			UserName: c.UserName,
		}
		hub.Broadcast <- msg
	}
}
