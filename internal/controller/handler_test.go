package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewHandler - тест для функции NewHandler
func TestNewHandler(t *testing.T) {
	mockService := new(MockAggregationService)

	handler := NewHandler(mockService)

	assert.NotNil(t, handler)
	assert.Equal(t, mockService, handler.aggregationService)
}
