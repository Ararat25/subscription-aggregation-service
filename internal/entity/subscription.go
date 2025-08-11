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
	Id          int       `json:"id,omitempty" example:"1" validate:"omitempty"`                                    // id подписки в бд
	ServiceName string    `json:"service_name"  example:"Netflix"  validate:"required"`                             // название сервиса, предоставляющего подписку
	Price       int       `json:"price" example:"499" validate:"required,min=1"`                                    // стоимость месячной подписки в рублях
	UserId      uuid.UUID `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" validate:"required,uuid4"` // id пользователя в формате UUID
	StartDate   string    `json:"start_date" example:"08-2025" validate:"required,datetime=01-2006"`                // дата начала подписки (месяц и год)
	EndDate     *string   `json:"end_date,omitempty" example:"09-2025" validate:"omitempty,datetime=01-2006"`       // дата окончания подписки (месяц и год)
}
