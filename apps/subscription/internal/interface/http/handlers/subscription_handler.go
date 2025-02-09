package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/subscription-tracker/subscription/internal/core/application"
	"github.com/subscription-tracker/subscription/internal/interface/http/dto"
)

type SubscriptionHandler struct {
	service *application.SubscriptionService
}

func NewSubscriptionHandler(service *application.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: service}
}

func (h *SubscriptionHandler) GetSubscription(c *gin.Context) {
	uuid := c.Param("uuid")
	userId := c.GetString("user_id")
	subscription, err := h.service.GetSubscription(uuid, userId)
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to fetch subscription"})
		return
	}
	c.JSON(200, dto.FromSubscription(subscription))
}

func (h *SubscriptionHandler) GetSubscriptions(c *gin.Context) {
	userId, _ := c.Get("user_id")
	var params dto.SubscriptionQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(400, gin.H{"error": "Invalid query parameters"})
		return
	}

	subscriptions, err := h.service.GetUserSubscriptions(params, userId.(string))
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch subscriptions"})
		return
	}

	c.JSON(200, subscriptions)
}

func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	var request dto.CreateSubscriptionRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	userID := c.GetString("user_id")
	subscription := request.ToSubscription(userID)

	if err := h.service.CreateSubscription(subscription); err != nil {
		c.JSON(500, gin.H{"error": "Failed to create subscription"})
		return
	}
	c.JSON(201, dto.FromSubscription(subscription))
}

func (h *SubscriptionHandler) UpdateSubscription(c *gin.Context) {
	var request dto.UpdateSubscriptionRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}
	id := c.Param("uuid")
	userId := c.GetString("user_id")
	err := h.service.UpdateSubscription(id, request, userId)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update subscription"})
		return
	}
	subscription, err := h.service.GetSubscription(id, userId)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch subscription"})
		return
	}
	c.JSON(200, dto.FromSubscription(subscription))
}

func (h *SubscriptionHandler) DeleteSubscription(c *gin.Context) {
	id := c.Param("id")
	userId := c.GetString("user_id")
	if err := h.service.DeleteSubscription(id, userId); err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete subscription"})
		return
	}
	c.Status(204)
}
