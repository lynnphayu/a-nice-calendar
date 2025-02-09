package application

import (
	"errors"
	"fmt"
	"sort"

	"github.com/google/uuid"
	"github.com/subscription-tracker/subscription/internal/core/domain"
	"github.com/subscription-tracker/subscription/internal/interface/http/dto"
)

type SubscriptionService struct {
	repo domain.SubscriptionRepository
}

func NewSubscriptionService(repo domain.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

func (s *SubscriptionService) GetSubscription(uuid string, userId string) (*domain.Subscription, error) {
	order := "created_at DESC"
	subs, err := s.repo.Find(&domain.SubscriptionRepoQuery{
		Uuid:   &uuid,
		UserID: &userId,
	}, &order)
	if err != nil {
		return nil, err
	}
	return &subs[0], nil
}

func (s *SubscriptionService) UpdateSubscription(uuid string, data dto.UpdateSubscriptionRequest, userId string) error {
	subs, err := s.repo.FindByUuid(uuid)
	if err != nil {
		return err
	}
	if subs.UserID != userId {
		return errors.New("Unauthorized")
	}
	// Create a new version of the subscription with a new ID but preserving the UUID
	newSubs := &domain.Subscription{
		ID:           0, // Let the database assign a new ID
		UserID:       userId,
		Uuid:         uuid,
		Name:         data.Name,
		Price:        data.Price,
		StartDate:    data.StartDate,
		BillingCycle: subs.BillingCycle,
		Logo:         data.Logo,
	}
	return s.repo.Create(newSubs)
}

func (s *SubscriptionService) GetUserSubscriptions(query dto.SubscriptionQueryParams, userId string) ([]domain.Subscription, error) {
	q := &domain.SubscriptionRepoQuery{
		UserID: &userId,
	}

	fmt.Println(query)
	if query.StartDateFrom != nil {
		q.StartDateFrom = query.StartDateFrom
	}
	if query.StartDateTo != nil {
		q.StartDateTo = query.StartDateTo
	}
	// if query.Day != nil && query.Month != nil && query.Year != nil {
	// 	to := time.Date(*query.Year, time.Month(*query.Month), *query.Day, 0, 0, 0, 0, time.UTC)
	// 	q.StartDateTo = &to
	// } else if query.Month != nil && query.Year != nil {
	// 	from := time.Date(*query.Year, time.Month(*query.Month), 1, 0, 0, 0, 0, time.UTC)
	// 	to := time.Date(*query.Year, time.Month(*query.Month+1), 1, 0, 0, 0, 0, time.UTC)
	// 	to = to.Add(-24 * time.Hour)
	// 	q.StartDateFrom = &from
	// 	q.StartDateTo = &to
	// } else if query.Year != nil {
	// 	from := time.Date(*query.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	// 	to := time.Date(*query.Year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	// 	to = to.Add(-24 * time.Hour)
	// 	q.StartDateFrom = &from
	// 	q.StartDateTo = &to
	// }
	fmt.Println(q)
	subs, err := s.repo.Find(q, nil)
	if err != nil {
		return nil, err
	}
	uuidGrouped := make(map[string][]domain.Subscription)
	for _, sub := range subs {
		uuidGrouped[sub.Uuid] = append(uuidGrouped[sub.Uuid], sub)
	}
	latestSubs := make([]domain.Subscription, 0)
	for _, subs := range uuidGrouped {
		sort.Slice(subs, func(i, j int) bool {
			return subs[i].StartDate.After(subs[j].StartDate)
		})
		latestSubs = append(latestSubs, subs[0])
	}

	return latestSubs, nil
}

func (s *SubscriptionService) CreateSubscription(subscription *domain.Subscription) error {
	if subscription.Uuid == "" {
		subscription.Uuid = uuid.NewString()
	}
	return s.repo.Create(subscription)
}

func (s *SubscriptionService) DeleteSubscription(id string, userId string) error {
	return s.repo.Delete(id)
}
