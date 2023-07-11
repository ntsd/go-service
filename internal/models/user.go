package models

import "time"

// User model
// @Description User information includes id, email, and name.
type User struct {
	ID        string    `gorm:"column:id"`
	Email     string    `gorm:"column:email"`
	Name      string    `gorm:"column:name"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (m *User) TableName() string {
	return "users"
}
