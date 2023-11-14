package Controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go-chat/Models"
)

func (h *Handler) SendMessageToQueue(c *gin.Context, m *Models.Message) {
	messageJSON, err := json.Marshal(m)
	if err != nil {
		fmt.Println("Error encoding message to JSON:", err)
	}
	option := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	_, errs := option.LPush(c, "chat-queue", string(messageJSON)).Result()
	if errs != nil {
		fmt.Println("Error pushing message to Redis queue:", err)
		return
	}
}

func (h *Handler) ReceiveMessageMq(c *gin.Context) {
	option := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	for {
		rcvmessage, err := option.BRPop(c, 0, "chat-queue").Result()
		if err != nil {
			panic(err)
		}
		h.HandleRedisMessage(rcvmessage)

	}

}
func (h *Handler) HandleRedisMessage(message []string) {
	// Convert message from slice of strings to slice of bytes
	var messageBytes [][]byte
	for _, s := range message {
		messageBytes = append(messageBytes, []byte(s))
	}

	// Decode message from JSON to a Models.Message struct
	var m Models.Message
	if err := json.Unmarshal(bytes.Join(messageBytes, []byte("")), &m); err != nil {
		fmt.Println("Error decoding Redis message:", err)
		return
	}

	// Handle the Redis message
	h.hub.Broadcast <- &m
}
