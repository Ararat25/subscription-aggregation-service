package model

import (
	"github.com/Ararat25/subscription-aggregation-service/internal/entity"
	"github.com/Ararat25/subscription-aggregation-service/internal/repository"
	"github.com/google/uuid"
	"time"
)

const dateLayout = "01-2006"

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

func (ags *AggregationService) CreateSubscription(s *entity.Subscription) error {
	err := ags.Storage.CreateSubscription(s)
	if err != nil {
		return err
	}

	return nil
}

func (ags *AggregationService) ReadSubscription(id int64) (*entity.Subscription, error) {
	sub, err := ags.Storage.ReadSubscription(id)
	if err != nil {
		return nil, err
	}

	return sub, nil
}

func (ags *AggregationService) UpdateSubscription(s *entity.Subscription) error {
	err := ags.Storage.UpdateSubscription(s)
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
