package postgres

import (
	"fmt"

	"github.com/subscription-tracker/subscription/internal/core/domain"
	"gorm.io/gorm"
)

type SubscriptionConfigRepository struct {
	db *gorm.DB
}

func NewSubscriptionConfigRepository(db *gorm.DB) *SubscriptionConfigRepository {
	return &SubscriptionConfigRepository{db: db}
}

func (r *SubscriptionConfigRepository) Find() (*[]domain.SubscriptionConfig, error) {
	var subsConfig []domain.SubscriptionConfig
	result := r.db.Model(&domain.SubscriptionConfig{}).Preload("Plans").Find(&subsConfig)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch subscription configs: %w", result.Error)
	}
	return &subsConfig, nil
}

func (r *SubscriptionConfigRepository) FindByProvider(provider string) (*domain.SubscriptionConfig, error) {
	var subsConfig domain.SubscriptionConfig
	result := r.db.Model(&domain.SubscriptionConfig{}).Preload("Plans").Where("provider = ?", provider).First(&subsConfig)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch subscription config: %w", result.Error)
	}
	return &subsConfig, nil
}

func (r *SubscriptionConfigRepository) FindActivePlans() (*[]domain.SubscriptionConfigPlan, error) {
	var plans []domain.SubscriptionConfigPlan
	result := r.db.Model(&domain.SubscriptionConfigPlan{}).
		Joins("JOIN subscription_configs ON subscription_configs.id = subscription_config_plans.subscription_config_id").
		Where("subscription_config_plans.status = ? AND subscription_configs.status = ?", "active", "active").
		Find(&plans)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch active plans: %w", result.Error)
	}
	return &plans, nil
}

func (r *SubscriptionConfigRepository) Create(config *domain.SubscriptionConfig) error {
	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid subscription config: %w", err)
	}

	result := r.db.Create(config)
	if result.Error != nil {
		return fmt.Errorf("failed to create subscription config: %w", result.Error)
	}
	return nil
}

func (r *SubscriptionConfigRepository) Update(config *domain.SubscriptionConfig) error {
	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid subscription config: %w", err)
	}

	result := r.db.Save(config)
	if result.Error != nil {
		return fmt.Errorf("failed to update subscription config: %w", result.Error)
	}
	return nil
}

func (r *SubscriptionConfigRepository) Delete(id int64) error {
	result := r.db.Delete(&domain.SubscriptionConfig{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete subscription config: %w", result.Error)
	}
	return nil
}
