package model

import (
	"context"
	"fmt"
	"time"

	"github.com/Ararat25/subscription-aggregation-service/internal/entity"
	"github.com/Ararat25/subscription-aggregation-service/internal/repository"
	"github.com/google/uuid"
)

const DateLayout = "01-2006" // шаблон для преобразования дат из строки в объект

// AggregationService - структура для сервиса агрегации
type AggregationService struct {
	Storage repository.Repo // объект для работы с бд
}

// NewAggregationService возвращает новый объект структуры Service
func NewAggregationService(storage repository.Repo) *AggregationService {
	return &AggregationService{
		Storage: storage,
	}
}

// CreateSubscription добавляет подписку в бд и возвращает id
func (ags *AggregationService) CreateSubscription(ctx context.Context, s *entity.SubscriptionRequest) (int64, error) {
	if s == nil {
		return 0, fmt.Errorf("invalid argument error")
	}

	subNew, err := convertStringDateToTime(s)
	if err != nil {
		return 0, err
	}

	if !isEndDateValid(subNew.StartDate, *subNew.EndDate) {
		return 0, fmt.Errorf("end_date must be >= start_date")
	}

	id, err := ags.Storage.CreateSubscription(ctx, subNew)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// ReadSubscription возвращает подписку из бд по id
func (ags *AggregationService) ReadSubscription(ctx context.Context, id int64) (*entity.Subscription, error) {
	sub, err := ags.Storage.ReadSubscription(ctx, id)
	if err != nil {
		return nil, err
	}

	return sub, nil
}

// UpdateSubscription обновляет данные подписки в бд
func (ags *AggregationService) UpdateSubscription(ctx context.Context, s *entity.SubscriptionRequest) error {
	if s == nil {
		return fmt.Errorf("invalid argument error")
	}

	subNew, err := convertStringDateToTime(s)
	if err != nil {
		return err
	}

	if !isEndDateValid(subNew.StartDate, *subNew.EndDate) {
		return fmt.Errorf("end_date must be >= start_date")
	}

	err = ags.Storage.UpdateSubscription(ctx, subNew)
	if err != nil {
		return err
	}

	return nil
}

// DeleteSubscription удаляет подписку из бд
func (ags *AggregationService) DeleteSubscription(ctx context.Context, id int64) error {
	err := ags.Storage.DeleteSubscription(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

// ListSubscriptions возвращает список всех подписок из бд
func (ags *AggregationService) ListSubscriptions(ctx context.Context) ([]*entity.Subscription, error) {
	subs, err := ags.Storage.ListSubscriptions(ctx)
	if err != nil {
		return nil, err
	}

	return subs, nil
}

// TotalCost возвращает суммарную стоимость подписок за определенный период с фильтрацией по id пользователя и названию сервиса
func (ags *AggregationService) TotalCost(ctx context.Context, from, to time.Time, userID *uuid.UUID, serviceName *string) (int, error) {
	fromReset := resetDay(from)
	toReset := resetDay(to)

	if !isEndDateValid(fromReset, toReset) {
		return 0, fmt.Errorf("to must be >= from")
	}

	cost, err := ags.Storage.TotalCost(ctx, fromReset, toReset, userID, serviceName)
	if err != nil {
		return 0, err
	}

	return cost, nil
}

// resetDay обнуляет день
func resetDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// convertStringDateToTime преобразует даты в подписке из строки в time.Time
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

// isEndDateValid проверяет, что endDate >= startDate
func isEndDateValid(startDate, endDate time.Time) bool {
	return !endDate.Before(startDate)
}
