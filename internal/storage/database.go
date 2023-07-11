package storage

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewDatabase returns gorm database connection
func NewDatabase(postgresURL string) (*gorm.DB, error) {
	var finalErr error

	// create retry to connect to database with delay 3 seconds
	for tires := 0; tires < 3; tires++ {
		db, err := gorm.Open(postgres.Open(postgresURL), &gorm.Config{
			// GORM performs write operations inside a transaction to ensure data consistency.
			SkipDefaultTransaction: true,
			// allow statement to cache
			PrepareStmt: true,
		})
		if err != nil {
			finalErr = err
			time.Sleep(3 * time.Second)
			continue
		}
		return db, nil
	}

	return nil, finalErr
}
