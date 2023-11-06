package Models

type CreateRoomReq struct {
	ID     string `gorm:"AutoIncrement" json:"id"`
	RoomID string `gorm:"foreignKey" json:"roomId"`
	Name   string `json:"name"`
}
