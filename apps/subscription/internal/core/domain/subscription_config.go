package domain

import (
	"errors"
	"time"
)

// SubscriptionConfig represents a subscription provider configuration
type SubscriptionConfig struct {
	ID          uint                     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Provider    string                   `gorm:"column:provider;not null;uniqueIndex" json:"provider" validate:"required"`
	Description string                   `gorm:"column:description" json:"description"`
	Logo        string                   `gorm:"column:logo" json:"logo" validate:"url"`
	Website     string                   `gorm:"column:website" json:"website" validate:"url"`
	Status      string                   `gorm:"column:status;default:'active'" json:"status" validate:"oneof=active inactive deprecated"`
	Plans       []SubscriptionConfigPlan `gorm:"foreignKey:SubscriptionConfigID" json:"plans"`

	CreatedAt time.Time `gorm:"column:created_at;not null;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;autoUpdateTime" json:"updatedAt"`
}

// SubscriptionConfigPlan represents a specific plan within a subscription configuration
type SubscriptionConfigPlan struct {
	ID                   uint    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name                 string  `gorm:"column:name;not null" json:"name" validate:"required"`
	Description          string  `gorm:"column:description" json:"description"`
	Price                float64 `gorm:"column:price;not null" json:"price" validate:"required,gte=0"`
	Currency             string  `gorm:"column:currency;not null" json:"currency" validate:"required,iso4217"`
	BillingCycle         int32   `gorm:"column:billing_cycle;not null" json:"billing_cycle" validate:"required,gte=1"`
	Status               string  `gorm:"column:status;default:'active'" json:"status" validate:"oneof=active inactive deprecated"`
	SubscriptionConfigID uint    `gorm:"column:subscription_config_id;not null" json:"subscription_config_id"`

	CreatedAt time.Time `gorm:"column:created_at;not null;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;autoUpdateTime" json:"updatedAt"`
}

// Validate performs basic validation on the subscription config
func (sc *SubscriptionConfig) Validate() error {
	if sc.Provider == "" {
		return errors.New("provider is required")
	}
	return nil
}

// IsActive checks if the subscription config is active
func (sc *SubscriptionConfig) IsActive() bool {
	return sc.Status == "active"
}

// GetActivePlans returns only the active plans from the subscription config
func (sc *SubscriptionConfig) GetActivePlans() []SubscriptionConfigPlan {
	var activePlans []SubscriptionConfigPlan
	for _, plan := range sc.Plans {
		if plan.Status == "active" {
			activePlans = append(activePlans, plan)
		}
	}
	return activePlans
}

// Validate performs basic validation on the subscription plan
func (scp *SubscriptionConfigPlan) Validate() error {
	if scp.Name == "" {
		return errors.New("name is required")
	}
	if scp.Price < 0 {
		return errors.New("price must be non-negative")
	}
	if scp.BillingCycle < 1 {
		return errors.New("billing cycle must be positive")
	}
	return nil
}

// IsActive checks if the plan is active
func (scp *SubscriptionConfigPlan) IsActive() bool {
	return scp.Status == "active"
}

type SubscriptionConfigRepository interface {
	Find() (*[]SubscriptionConfig, error)
	FindByProvider(provider string) (*SubscriptionConfig, error)
	FindActivePlans() (*[]SubscriptionConfigPlan, error)
	Create(config *SubscriptionConfig) error
	Update(config *SubscriptionConfig) error
	Delete(id int64) error
}
