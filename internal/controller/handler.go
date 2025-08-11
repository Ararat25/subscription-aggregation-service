package controller

import (
	"github.com/Ararat25/subscription-aggregation-service/internal/model"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// Handler структура для обработчиков запросов
type Handler struct {
	aggregationService *model.AggregationService // объект для работы с сервисом агрегации подписок
}

// NewHandler создает новый объект Handler
func NewHandler(aggregationService *model.AggregationService) *Handler {
	return &Handler{
		aggregationService: aggregationService,
	}
}
