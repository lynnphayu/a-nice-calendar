package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Config holds the database configuration
type Config struct {
	URL string
}

// Database represents the database connection
type Database struct {
	client *mongo.Client
}

// NewDatabase creates a new database connection
func NewDatabase(config Config) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("Database URL:", config.URL)
	// Create MongoDB client
	clientOptions := options.Client().ApplyURI(config.URL)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping the database
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return &Database{client: client}, nil
}

// GetDB returns the mongo client
func (d *Database) GetDB() *mongo.Client {
	return d.client
}

// Close closes the database connection
func (d *Database) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return d.client.Disconnect(ctx)
}
