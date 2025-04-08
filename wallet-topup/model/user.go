package model

type User struct {
	UserID  uint    `gorm:"primaryKey"`
	Balance float64 `gorm:"type:numeric(12,2)"`
}

type UserRepository interface {
	GetUserByID(userID uint) (*User, error)
	UpdateUserBalance(userID uint, amount float64) error
}
