package storage

import (
	"go-service/internal/models"

	"gorm.io/hints"
)

// CreateUser create User to the database
func (s *storage) CreateUser(user *models.User) error {
	return s.db.Create(user).Error
}

// ListUsers list all user from the database
func (s *storage) ListUsers(users *[]models.User, offset, limit int, name string) error {
	if name != "" {
		return s.db.Clauses(hints.UseIndex("users_name_trgm_idx")).
			Where("WORD_SIMILARITY(?, name) > 0.4", name).
			Offset(offset).
			Limit(limit).
			Find(users).
			Error
	}
	return s.db.Offset(offset).Limit(limit).Find(users).Error
}

// GetUserByID get one User by id
func (s *storage) GetUserByID(user *models.User, id string) error {
	return s.db.First(user, "id = ?", id).Error
}

// UpdateUserByID update one User by id
func (s *storage) UpdateUserByID(user *models.User, id string) error {
	return s.db.Where("id = ?", id).Updates(user).First(user).Error
}
