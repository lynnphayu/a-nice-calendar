package main

import (
	"log"
	"os"

	_ "ariga.io/atlas-provider-gorm/gormschema"
	"github.com/joho/godotenv"
	"github.com/subscription-tracker/subscription/internal/app"
	postgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Initialize database connection
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run database migrations
	// if err := migrations.RunMigrations(db); err != nil {
	// 	log.Fatalf("Failed to run migrations: %v", err)
	// }

	// Create and configure application
	application := app.NewApplication(db)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := application.Router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
