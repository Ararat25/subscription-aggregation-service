package entity

import (
	"time"

	"github.com/google/uuid"
)

// Subscription - структура для храненеия данных подписки
type Subscription struct {
	Id          int        `json:"-"`                  // id подписки в бд
	ServiceName string     `json:"service_name"`       // название сервиса, предоставляющего подписку
	Price       int        `json:"price"`              // стоимость месячной подписки в рублях
	UserId      uuid.UUID  `json:"user_id"`            // id пользователя в формате UUID
	StartDate   time.Time  `json:"start_date"`         // дата начала подписки (месяц и год)
	EndDate     *time.Time `json:"end_date,omitempty"` // дата окончания подписки (месяц и год)
}

// SubscriptionRequest - структура для парсинга данных подписки из запроса
type SubscriptionRequest struct {
	Id          int       `json:"id,omitempty"`       // id подписки в бд
	ServiceName string    `json:"service_name"`       // название сервиса, предоставляющего подписку
	Price       int       `json:"price"`              // стоимость месячной подписки в рублях
	UserId      uuid.UUID `json:"user_id"`            // id пользователя в формате UUID
	StartDate   string    `json:"start_date"`         // дата начала подписки (месяц и год)
	EndDate     *string   `json:"end_date,omitempty"` // дата окончания подписки (месяц и год)
}
