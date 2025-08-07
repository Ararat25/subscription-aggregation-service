package model

import (
	"gorm.io/gorm"
)

// AggregationService - структура для сервиса агрегации
type AggregationService struct {
	Storage *gorm.DB
}

// NewAggregationService возвращает новый объект структуры Service
func NewAggregationService(storage *gorm.DB) *AggregationService {
	return &AggregationService{
		Storage: storage,
	}
}
