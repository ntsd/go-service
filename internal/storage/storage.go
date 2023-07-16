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
	writeDB *gorm.DB
	readDB  *gorm.DB
}

// NewStorage create a new storage
func NewStorage(cfg *config.Config) (Storage, error) {
	writeDB, err := NewDatabase(cfg.PostgresWriteURL)
	if err != nil {
		return nil, fmt.Errorf("error to create write database: %w", err)
	}
	if cfg.DevMode {
		writeDB = writeDB.Debug()
	}

	var readDB *gorm.DB
	if cfg.PostgresWriteURL == cfg.PostgresReadURL {
		readDB = writeDB
	} else {
		db, err := NewDatabase(cfg.PostgresReadURL)
		if err != nil {
			return nil, fmt.Errorf("error to create read database: %w", err)
		}
		if cfg.DevMode {
			writeDB = writeDB.Debug()
		}
		readDB = db
	}

	return &storage{
		writeDB: writeDB,
		readDB:  readDB,
	}, nil
}
