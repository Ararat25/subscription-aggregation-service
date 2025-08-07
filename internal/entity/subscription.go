package entity

import (
	"github.com/google/uuid"
	"time"
)

// Subscription - структура для храненеия данных подписки
type Subscription struct {
	id          int       `gorm:"column:id;primaryKey"`
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserId      uuid.UUID `json:"user_id"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date,omitempty"`
}
