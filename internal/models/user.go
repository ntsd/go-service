package models

import "time"

// User model
// @Description User information includes id, email, and name.
type User struct {
	ID        string    `gorm:"column:id" json:"id"`
	Email     string    `gorm:"column:email" json:"email"`
	Name      string    `gorm:"column:name" json:"name"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (m *User) TableName() string {
	return "users"
}
