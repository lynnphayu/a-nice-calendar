package dto

import (
	"time"

	"github.com/subscription-tracker/subscription/internal/core/domain"
)

type CreateSubscriptionRequest struct {
	Name         string    `json:"name" binding:"required"`
	Price        float64   `json:"price" binding:"required,gte=0"`
	BillingCycle int       `json:"billingCycle" binding:"required,gte=0"`
	StartDate    time.Time `json:"startDate" binding:"required"`
	Logo         string    `json:"logo" binding:"required"`
}

type UpdateSubscriptionRequest struct {
	Name      string    `json:"name" binding:"required"`
	Price     float64   `json:"price" binding:"required,gte=0"`
	StartDate time.Time `json:"startDate" binding:"required"`
	Logo      string    `json:"logo" binding:"required"`
}

type SubscriptionQueryParams struct {
	StartDateFrom *time.Time `form:"start_date_from" time_format:"2006-01-02"`
	StartDateTo   *time.Time `form:"start_date_to" time_format:"2006-01-02" `
}

type SubscriptionResponse struct {
	UUID         string    `json:"uuid"`
	Name         string    `json:"name"`
	Price        float64   `json:"price"`
	BillingCycle int       `json:"billingCycle"`
	StartDate    time.Time `json:"startDate"`
	Logo         string    `json:"logo"`
	UserID       string    `json:"userId"`
}

// ToSubscription converts CreateSubscriptionRequest to domain.Subscription
func (r *CreateSubscriptionRequest) ToSubscription(userID string) *domain.Subscription {
	return &domain.Subscription{
		Name:         r.Name,
		Price:        r.Price,
		BillingCycle: r.BillingCycle,
		StartDate:    r.StartDate,
		Logo:         r.Logo,
		UserID:       userID,
	}
}

// FromSubscription creates SubscriptionResponse from domain.Subscription
func FromSubscription(s *domain.Subscription) *SubscriptionResponse {
	return &SubscriptionResponse{
		UUID:         s.Uuid,
		Name:         s.Name,
		Price:        s.Price,
		BillingCycle: s.BillingCycle,
		StartDate:    s.StartDate,
		Logo:         s.Logo,
		UserID:       s.UserID,
	}
}
