package Databases

import (
	"go-chat/Models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var CDB *gorm.DB

func ConnectClientUser() {
	connectClient, err := gorm.Open(mysql.Open("root:@tcp(localhost:3306)/chat-realtime"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	CDB = connectClient
	connectClient.AutoMigrate(&Models.UserClients{})
}
