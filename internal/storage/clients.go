package storage

import (
	"go-service/internal/models"
)

// CreateClient create Client to the database
func (s *storage) CreateClient(client *models.Client) error {
	return s.db.Create(client).Error
}

// GetClientByID get one Client by id
func (s *storage) GetClientByID(client *models.Client, id string) error {
	return s.db.First(client, "id = ?", id).Error
}
