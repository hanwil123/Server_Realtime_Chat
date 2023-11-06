package Models

type UserClients struct {
	ID       uint64 `gorm:"primaryKey;AutoIncrement" json:"id"`
	UserID   uint64 `gorm:"foreignKey" json:"userId"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password []byte
}
