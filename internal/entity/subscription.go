package entity

import (
	"time"

	"github.com/google/uuid"
)

// Subscription - структура для храненеия данных подписки
type Subscription struct {
	Id          int        `json:"-"`
	ServiceName string     `json:"service_name"`
	Price       int        `json:"price"`
	UserId      uuid.UUID  `json:"user_id"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}

type SubscriptionRequest struct {
	Id          int       `json:"id,omitempty"`
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserId      uuid.UUID `json:"user_id"`
	StartDate   string    `json:"start_date"`
	EndDate     *string   `json:"end_date,omitempty"`
}
