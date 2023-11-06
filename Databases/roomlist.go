package Databases

import (
	"go-chat/Models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var RDB *gorm.DB

func ConnectRoomList() {
	connectRoom, err := gorm.Open(mysql.Open("root:@tcp(localhost:3306)/chat-realtime"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	RDB = connectRoom
	connectRoom.AutoMigrate(&Models.CreateRoomReq{})
}
