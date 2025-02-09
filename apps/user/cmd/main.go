package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/subscription-tracker/user/internal/app"
	"github.com/subscription-tracker/user/internal/database"
)

func main() {
	// Initialize database configuration

	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	dbConfig := database.Config{
		URL: os.Getenv("DATABASE_URL"),
	}
	fmt.Println(dbConfig.URL)

	// Create new database instance
	db, err := database.NewDatabase(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create new application instance
	application, err := app.NewApplication(db.GetDB())
	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := application.Router.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
