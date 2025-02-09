package domain

import "time"

type Subscription struct {
	ID           uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Uuid         string    `gorm:"column:uuid; not null" json:"uuid"`
	Name         string    `gorm:"column:name;not null" json:"name"`
	Price        float64   `gorm:"column:price;not null" json:"price"`
	BillingCycle int       `gorm:"column:billing_cycle;not null" json:"billingCycle"`
	StartDate    time.Time `gorm:"column:start_date;not null" json:"startDate"`
	Logo         string    `gorm:"column:logo" json:"logo"`
	UserID       string    `gorm:"column:user_id;not null" json:"userId"`
	IsActive     bool      `gorm:"column:is_active;not null;default:true" json:"isActive"`

	CreatedAt time.Time `gorm:"column:created_at;not null;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;autoUpdateTime" json:"updatedAt"`
}

type SubscriptionRepoQuery struct {
	StartDateFrom *time.Time
	StartDateTo   *time.Time
	UserID        *string
	BillingCycle  *string
	Name          *string
	Uuid          *string
}

type SubscriptionRepository interface {
	Find(query *SubscriptionRepoQuery, order *string) ([]Subscription, error)
	FindByUserId(userId string) ([]Subscription, error)
	FindByUuid(uuid string) (*Subscription, error)
	Create(subscription *Subscription) error
	Delete(id string) error
}
