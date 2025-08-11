package error

import "errors"

var (
	ErrSubscriptionNotFound = errors.New("subscription not found") // подписка не найдена
)
