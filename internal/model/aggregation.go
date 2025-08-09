package model

import (
	"fmt"
	"time"

	"github.com/Ararat25/subscription-aggregation-service/internal/entity"
	"github.com/Ararat25/subscription-aggregation-service/internal/repository"
	"github.com/google/uuid"
)

const DateLayout = "01-2006"

// AggregationService - структура для сервиса агрегации
type AggregationService struct {
	Storage repository.Repo
}

// NewAggregationService возвращает новый объект структуры Service
func NewAggregationService(storage repository.Repo) *AggregationService {
	return &AggregationService{
		Storage: storage,
	}
}

func (ags *AggregationService) CreateSubscription(s *entity.SubscriptionRequest) (int64, error) {
	if s == nil {
		return 0, fmt.Errorf("invalid argument error")
	}

	subNew, err := convertStringDateToTime(s)
	if err != nil {
		return 0, err
	}

	id, err := ags.Storage.CreateSubscription(subNew)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (ags *AggregationService) ReadSubscription(id int64) (*entity.Subscription, error) {
	sub, err := ags.Storage.ReadSubscription(id)
	if err != nil {
		return nil, err
	}

	return sub, nil
}

func (ags *AggregationService) UpdateSubscription(s *entity.SubscriptionRequest) error {
	if s == nil {
		return fmt.Errorf("invalid argument error")
	}

	subNew, err := convertStringDateToTime(s)
	if err != nil {
		return err
	}

	err = ags.Storage.UpdateSubscription(subNew)
	if err != nil {
		return err
	}

	return nil
}

func (ags *AggregationService) DeleteSubscription(id int64) error {
	err := ags.Storage.DeleteSubscription(id)
	if err != nil {
		return err
	}

	return nil
}

func (ags *AggregationService) ListSubscriptions() ([]*entity.Subscription, error) {
	subs, err := ags.Storage.ListSubscriptions()
	if err != nil {
		return nil, err
	}

	return subs, nil
}

func (ags *AggregationService) TotalCost(from, to time.Time, userID *uuid.UUID, serviceName *string) (int, error) {
	cost, err := ags.Storage.TotalCost(resetDay(from), resetDay(to), userID, serviceName)
	if err != nil {
		return 0, err
	}

	return cost, nil
}

func resetDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

func convertStringDateToTime(s *entity.SubscriptionRequest) (*entity.Subscription, error) {
	subNew := &entity.Subscription{
		ServiceName: s.ServiceName,
		Price:       s.Price,
		UserId:      s.UserId,
	}

	if s.Id != 0 {
		subNew.Id = s.Id
	}

	start, err := time.Parse(DateLayout, s.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid date format")
	}

	var end *time.Time
	if s.EndDate != nil {
		endValue, err := time.Parse(DateLayout, *s.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid date format")
		}

		end = &endValue
	}

	subNew.StartDate = start
	subNew.EndDate = end

	return subNew, nil
}
