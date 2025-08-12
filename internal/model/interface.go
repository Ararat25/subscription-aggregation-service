package model

import (
	"context"
	"time"

	"github.com/Ararat25/subscription-aggregation-service/internal/entity"
	"github.com/google/uuid"
)

// Service интерфейс для сервиса агрегации
type Service interface {
	CreateSubscription(ctx context.Context, s *entity.SubscriptionRequest) (int64, error)
	ReadSubscription(ctx context.Context, id int64) (*entity.Subscription, error)
	UpdateSubscription(ctx context.Context, s *entity.SubscriptionRequest) error
	DeleteSubscription(ctx context.Context, id int64) error
	ListSubscriptions(ctx context.Context) ([]*entity.Subscription, error)
	TotalCost(ctx context.Context, from, to time.Time, userID *uuid.UUID, serviceName *string) (int, error)
}
