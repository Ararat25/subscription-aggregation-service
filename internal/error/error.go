package error

import (
	"errors"
)

var (
	ErrSubscriptionNotFound = errors.New("subscription not found")         // подписка не найдена
	ErrDateRange            = errors.New("end_date must be >= start_date") // дата конца должна быть >= дате начала
)
