package application

import "github.com/subscription-tracker/subscription/internal/core/domain"

type SubscriptionConfigService struct {
	repo domain.SubscriptionConfigRepository
}

func NewSubscriptionConfigService(repo domain.SubscriptionConfigRepository) *SubscriptionConfigService {
	return &SubscriptionConfigService{repo: repo}
}

func (s *SubscriptionConfigService) GetSubscriptionConfigs() (*[]domain.SubscriptionConfig, error) {
	subs, err := s.repo.Find()
	return subs, err
}
