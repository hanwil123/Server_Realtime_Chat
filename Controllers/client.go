package Controllers

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var UDB *gorm.DB

func ConnectUsers() {
	connectUser, err := gorm.Open(mysql.Open("root:@tcp(localhost:3306)/chat-realtime"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	UDB = connectUser
	connectUser.AutoMigrate(&Client{})
}
