package repository

import (
	"time"

	"github.com/Ararat25/subscription-aggregation-service/internal/entity"
	"github.com/google/uuid"
)

// Repo интерфейс хранилища для подписок
type Repo interface {
	ConnectDB(dbHost string, dbUser string, dbPassword string, dbName string, dbPort int) error
	CreateSubscription(s *entity.Subscription) (int64, error)
	ReadSubscription(id int64) (*entity.Subscription, error)
	UpdateSubscription(s *entity.Subscription) error
	DeleteSubscription(id int64) error
	ListSubscriptions() ([]*entity.Subscription, error)
	TotalCost(from, to time.Time, userID *uuid.UUID, serviceName *string) (int, error)
}
