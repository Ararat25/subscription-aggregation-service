package controller

import (
	"github.com/Ararat25/subscription-aggregation-service/internal/model"
)

// Handler структура для обработчиков запросов
type Handler struct {
	aggregationService *model.AggregationService
}

// NewHandler создает новый объект Handler
func NewHandler(aggregationService *model.AggregationService) *Handler {
	return &Handler{
		aggregationService: aggregationService,
	}
}
