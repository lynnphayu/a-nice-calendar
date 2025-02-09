package database

import (
	"log"
	"os"

	"github.com/subscription-tracker/subscription/internal/core/domain"
	postgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Config holds database configuration
type Config struct {
	URL string
}

// Database represents a database instance
type Database struct {
	db *gorm.DB
}

// NewDatabase creates a new database instance with the given configuration
func NewDatabase(config Config) (*Database, error) {
	dbURL := config.URL
	if dbURL == "" {
		dbURL = os.Getenv("DATABASE_URL")
		if dbURL == "" {
			dbURL = "host=localhost user=postgres password=postgres dbname=subscriptions port=5432 sslmode=disable"
		}
	}

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate the schema
	if err := db.AutoMigrate(&domain.Subscription{}); err != nil {
		log.Printf("Failed to migrate database: %v", err)
		return nil, err
	}

	return &Database{db: db}, nil
}

// GetDB returns the underlying gorm.DB instance
func (d *Database) GetDB() *gorm.DB {
	return d.db
}
