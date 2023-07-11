package storage

import (
	"fmt"
	"go-service/internal/config"
	"go-service/internal/models"

	"gorm.io/gorm"
)

// Storage is the storage interface
type Storage interface {
	CreateClient(client *models.Client) error
	GetClientByID(client *models.Client, id string) error

	CreateUser(user *models.User) error
	ListUsers(users *[]models.User, offset, limit int, name string) error
	GetUserByID(user *models.User, id string) error
	UpdateUserByID(user *models.User, id string) error
}

// storage is the storage struct
type storage struct {
	db *gorm.DB
}

// NewStorage create a new storage
func NewStorage(cfg *config.Config) (Storage, error) {
	db, err := NewDatabase(cfg.PostgresURL)
	if err != nil {
		return nil, fmt.Errorf("error to create database: %w", err)
	}

	// debug mode
	if cfg.DevMode {
		db = db.Debug()
	}

	return &storage{
		db: db,
	}, nil
}
