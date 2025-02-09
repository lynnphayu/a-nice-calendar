package app

import (
	"github.com/gin-gonic/gin"
	"github.com/subscription-tracker/subscription/internal/core/application"
	"github.com/subscription-tracker/subscription/internal/infrastructure/postgres"
	"github.com/subscription-tracker/subscription/internal/interface/http/handlers"
	"github.com/subscription-tracker/subscription/internal/middleware"
	"gorm.io/gorm"
)

// Application represents the main application container
type Application struct {
	Router *gin.Engine
	DB     *gorm.DB
}

// NewApplication creates and configures a new application instance
func NewApplication(db *gorm.DB) *Application {
	app := &Application{
		Router: gin.Default(),
		DB:     db,
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

	subscriptionRepo := postgres.NewSubscriptionRepository(db)
	subscriptionService := application.NewSubscriptionService(subscriptionRepo)
	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionService)

	subscriptionConfigRepo := postgres.NewSubscriptionConfigRepository(db)
	subscriptionConfigService := application.NewSubscriptionConfigService(subscriptionConfigRepo)
	subscriptionConfigHandler := handlers.NewSubscriptionConfigHandler(subscriptionConfigService)

	// Register routes
	api := app.Router.Group("/api")
	{
		// Protected subscription routes
		subscriptions := api.Group("/subscriptions", middleware.AuthMiddleware())
		{
			subscriptions.GET("/:uuid", subscriptionHandler.GetSubscription)
			subscriptions.GET("", subscriptionHandler.GetSubscriptions)
			subscriptions.POST("", subscriptionHandler.CreateSubscription)
			subscriptions.PUT("/:uuid", subscriptionHandler.UpdateSubscription)
			subscriptions.DELETE("/:uuid", subscriptionHandler.DeleteSubscription)
		}
		subscriptionConfigs := api.Group("/subscriptions_configs")
		{
			subscriptionConfigs.GET("", subscriptionConfigHandler.GetSubscriptionConfigs)
		}

		// Add more domain routes here as needed
	}

	return app
}
