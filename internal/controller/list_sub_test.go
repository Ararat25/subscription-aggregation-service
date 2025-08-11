package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ararat25/subscription-aggregation-service/internal/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestListSubscriptions - тест для функции ListSubscriptions контроллера
func TestListSubscriptions(t *testing.T) {
	mockService := new(MockAggregationService)

	handler := NewHandler(mockService)

	// Тестовый случай 1: Успешное получение списка подписок
	{
		expectedSubs := []*entity.Subscription{
			{
				Id:          1,
				ServiceName: "Service 1",
				Price:       100,
				UserId:      uuid.New(),
			},
			{
				Id:          2,
				ServiceName: "Service 2",
				Price:       200,
				UserId:      uuid.New(),
			},
		}

		req := httptest.NewRequest("GET", "/subscriptions", nil)
		rw := httptest.NewRecorder()

		mockService.On("ListSubscriptions", mock.Anything).Return(expectedSubs, nil).Once()

		handler.ListSubscriptions(rw, req)

		assert.Equal(t, http.StatusOK, rw.Code)

		var actualSubs []*entity.SubscriptionRequest
		_ = json.Unmarshal(rw.Body.Bytes(), &actualSubs)

		expected := make([]*entity.SubscriptionRequest, 0, len(expectedSubs))
		for _, sub := range expectedSubs {
			expected = append(expected, entity.ParseSubscriptionToRequest(sub))
		}

		assert.Equal(t, expected, actualSubs)
		mockService.AssertExpectations(t)
	}

	// Тестовый случай 2: Ошибка сервиса агрегации при получении списка
	{
		req := httptest.NewRequest("GET", "/subscriptions", nil)
		rw := httptest.NewRecorder()

		mockService.On("ListSubscriptions", mock.Anything).Return([]*entity.Subscription{}, errors.New("internal service error")).Once()

		handler.ListSubscriptions(rw, req)

		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		var errResp ErrorResponse
		_ = json.Unmarshal(rw.Body.Bytes(), &errResp)
		assert.Contains(t, errResp.Error, "internal service error")
		mockService.AssertExpectations(t)
	}
}
