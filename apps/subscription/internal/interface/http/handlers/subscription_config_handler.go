package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/subscription-tracker/subscription/internal/core/application"
)

type SubscriptionConfigHandler struct {
	subscriptionConfigService *application.SubscriptionConfigService
}

func NewSubscriptionConfigHandler(service *application.SubscriptionConfigService) *SubscriptionConfigHandler {
	return &SubscriptionConfigHandler{
		subscriptionConfigService: service,
	}
}

func (h *SubscriptionConfigHandler) GetSubscriptionConfigs(c *gin.Context) {
	subscriptionConfigs, err := h.subscriptionConfigService.GetSubscriptionConfigs()
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to fetch subscription"})
		return
	}
	c.JSON(200, subscriptionConfigs)
}
