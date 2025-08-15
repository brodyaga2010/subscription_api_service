package service

import (
	"context"
	"eff_mobile/internal/model"
	"eff_mobile/internal/repository"
	"time"
)

type SubscriptionService struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionService(repo repository.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

func (s *SubscriptionService) CreateSubscription(ctx context.Context, req *model.SubscriptionRequest) (*model.Subscription, error) {
	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		return nil, model.ErrDateFormat
	}

	var endDate *time.Time
	if req.EndDate != nil && *req.EndDate != "" {
		et, err := time.Parse("01-2006", *req.EndDate)
		if err != nil {
			return nil, model.ErrDateFormat
		}
		endDate = &et
	}

	sub := &model.Subscription{
		Service:   req.Service,
		Price:     req.Price,
		UserID:    req.UserID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	if err := s.repo.Create(ctx, sub); err != nil {
		return nil, err
	}

	return sub, nil
}

func (s *SubscriptionService) GetSubscription(ctx context.Context, id int) (*model.Subscription, error) {
	return s.repo.Get(ctx, id)
}

func (s *SubscriptionService) UpdateSubscription(ctx context.Context, id int, req *model.SubscriptionRequest) (*model.Subscription, error) {
	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		return nil, model.ErrDateFormat
	}

	var endDate *time.Time
	if req.EndDate != nil && *req.EndDate != "" {
		et, err := time.Parse("01-2006", *req.EndDate)
		if err != nil {
			return nil, model.ErrDateFormat
		}
		endDate = &et
	}

	sub := &model.Subscription{
		ID:        id,
		Service:   req.Service,
		Price:     req.Price,
		UserID:    req.UserID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	if err := s.repo.Update(ctx, sub); err != nil {
		return nil, err
	}
	return sub, nil
}

func (s *SubscriptionService) DeleteSubscription(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *SubscriptionService) ListSubscriptions(ctx context.Context) ([]model.Subscription, error) {
	return s.repo.List(ctx)
}

func (s *SubscriptionService) CalculateAmount(ctx context.Context, req model.SumRequest) (*model.SumResponse, error) {
	if req.From == "" || req.To == "" {
		return nil, model.ErrDateIsNull
	}

	fromTime, err := time.Parse("01-2006", req.From)
	if err != nil {
		return nil, model.ErrDateFormat
	}
	req.FromTime = fromTime

	toTime, err := time.Parse("01-2006", req.To)
	if err != nil {
		return nil, model.ErrDateFormat
	}
	req.ToTime = toTime

	if req.ToTime.Before(req.FromTime) {
		return nil, model.ErrDateBefore
	}

	return s.repo.CalculateAmount(ctx, &req)
}
