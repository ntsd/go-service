package main

import (
	"embed"
	"errors"
	"flag"
	"fmt"
	"go-service/internal/storage"
	"log"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"gorm.io/gorm"
)

// mode up for migrate up, down for migrate down
var mode = flag.String("mode", "up", "Example `-mode down`. `down` for migrate down, `up` for migrate up")

//go:embed migrations/*.sql
var migrations embed.FS

// logger logger for go migrate
type logger struct{}

func (l *logger) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func (l *logger) Verbose() bool {
	return false
}

// initMigration initializes the migration
func initMigration(gormDB *gorm.DB) (*migrate.Migrate, error) {
	source, err := httpfs.New(http.FS(migrations), "migrations")
	if err != nil {
		return nil, err
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, err
	}

	dbDriver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	return migrate.NewWithInstance("httpfs", source, "postgres", dbDriver)
}

func main() {
	flag.Parse()
	if mode == nil {
		log.Fatal("missing mode")
	}
	isUp := true
	switch *mode {
	case "up":
		isUp = true
	case "down":
		isUp = false
	default:
		log.Fatalf("mode must be `up` or `down` got %s", *mode)
	}

	fmt.Printf("Migrating %s\n", *mode)

	postgresURL := os.Getenv("POSTGRES_URL")
	if postgresURL == "" {
		log.Fatal(fmt.Errorf("`%s` env is required", "POSTGRES_URL"))
	}

	gormDb, err := storage.NewDatabase(postgresURL)
	if err != nil {
		log.Fatalf("error creating database: %v", err)
	}

	migration, err := initMigration(gormDb)
	if err != nil {
		log.Fatalf("error creating migration: %v", err)
	}

	migration.Log = &logger{}

	if isUp {
		if err := migration.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				log.Println(err)
				return
			}
			log.Fatal(err)
		}
	} else {
		if err := migration.Down(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				log.Println(err)
				return
			}
			log.Fatal(err)
		}
	}

	fmt.Printf("Migrated %s success\n", *mode)
}
