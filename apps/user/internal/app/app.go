package app

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/subscription-tracker/user/internal/core/application"
	"github.com/subscription-tracker/user/internal/infrastructure/email"
	"github.com/subscription-tracker/user/internal/infrastructure/mongodb"
	"github.com/subscription-tracker/user/internal/interface/http/handlers"
	"github.com/subscription-tracker/user/internal/interface/http/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

// Application represents the main application container
type Application struct {
	Router *gin.Engine
	DB     *mongo.Database
}

// NewApplication creates and configures a new application instance
func NewApplication(db *mongo.Client) (*Application, error) {

	app := &Application{
		Router: gin.Default(),
		DB:     db.Database("subs"),
	}

	// Configure CORS
	app.Router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Initialize repositories
	userRepo := mongodb.NewUserRepository(app.DB)
	emailService := email.NewEmailService(&email.SMTPConfig{
		Host:     "smtp.gmail.com",
		Port:     587,
		Username: "lynnphayu@gmail.com",
		Password: "blxz bjgn nvgp cbqf",
		From:     "subs@tracker.com",
	})

	// Initialize services
	userService := application.NewUserService(userRepo, emailService)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)

	// Register routes
	api := app.Router.Group("/api")
	{
		protected := api.Group("", middleware.AuthMiddleware([]byte(os.Getenv("JWT_SECRET"))))
		{
			protected.GET("/me", userHandler.GetCurrentUser)
		}
		// User routes
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)
		api.POST("/passwordless/initiate", userHandler.InitiatePasswordlessLogin)
		api.POST("/passwordless/verify", userHandler.VerifyLoginToken)
		api.PUT("/users/:id", userHandler.UpdateUser)
		api.DELETE("/users/:id", userHandler.DeleteUser)

		// Spotify OAuth routes
		api.GET("/auth/spotify", userHandler.InitiateLogin)
		api.GET("/spotify-callback", userHandler.HandleCallback)

	}

	return app, nil
}
