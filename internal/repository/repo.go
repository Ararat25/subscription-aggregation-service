package repository

import (
	"context"
	"time"

	"github.com/Ararat25/subscription-aggregation-service/internal/entity"
	"github.com/google/uuid"
)

// Repo интерфейс хранилища для подписок
type Repo interface {
	ConnectDB(ctx context.Context, dbHost string, dbUser string, dbPassword string, dbName string, dbPort int) error
	CreateSubscription(ctx context.Context, s *entity.Subscription) (int64, error)
	ReadSubscription(ctx context.Context, id int64) (*entity.Subscription, error)
	UpdateSubscription(ctx context.Context, s *entity.Subscription) error
	DeleteSubscription(ctx context.Context, id int64) error
	ListSubscriptions(ctx context.Context) ([]*entity.Subscription, error)
	TotalCost(ctx context.Context, from, to time.Time, userID *uuid.UUID, serviceName *string) (int, error)
	Close(ctx context.Context) error
}
