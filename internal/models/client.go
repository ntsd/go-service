package models

import (
	"time"
)

// Client client information
type Client struct {
	ID        string    `gorm:"column:id"`
	Secret    []byte    `gorm:"column:secret"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (m *Client) TableName() string {
	return "clients"
}
