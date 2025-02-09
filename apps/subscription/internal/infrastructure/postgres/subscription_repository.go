package postgres

import (
	"fmt"

	"github.com/subscription-tracker/subscription/internal/core/domain"
	"gorm.io/gorm"
)

type SubscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	db.AutoMigrate(&domain.Subscription{})
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) FindByUserId(userId string) ([]domain.Subscription, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}

	var subscriptions []domain.Subscription
	result := r.db.Find(&subscriptions, "user_id = ?", userId)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch subscriptions: %w", result.Error)
	}

	return subscriptions, nil
}

func (r *SubscriptionRepository) Find(query *domain.SubscriptionRepoQuery, order *string) ([]domain.Subscription, error) {
	db := r.db

	if query != nil {
		if query.StartDateFrom != nil {
			db = db.Where("start_date >= ?", query.StartDateFrom)
		}
		if query.StartDateTo != nil {
			db = db.Where("start_date <= ?", query.StartDateTo)
		}
		if query.UserID != nil {
			db = db.Where("user_id = ?", *query.UserID)
		}
		if query.BillingCycle != nil {
			db = db.Where("billing_cycle = ?", *query.BillingCycle)
		}
		if query.Name != nil {
			db = db.Where("name like ?", "%"+*query.Name+"%")
		}
		if query.Uuid != nil {
			db = db.Where("uuid =?", *query.Uuid)
		}
	}

	var subscriptions []domain.Subscription

	result := db.Find(&subscriptions).Order(order)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch subscriptions: %w", result.Error)
	}

	return subscriptions, nil
}

func (r *SubscriptionRepository) Create(subscription *domain.Subscription) error {
	fmt.Println("subscription", subscription)
	return r.db.Create(subscription).Error
}

func (r *SubscriptionRepository) Delete(id string) error {
	return r.db.Delete(&domain.Subscription{}, "id = ?", id).Error
}

func (r *SubscriptionRepository) FindByUuid(uuid string) (*domain.Subscription, error) {
	var subscription domain.Subscription
	result := r.db.First(&subscription, "uuid = ?", uuid)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch subscription: %w", result.Error)
	}
	return &subscription, nil
}
