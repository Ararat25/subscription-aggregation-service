package controller

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Ararat25/subscription-aggregation-service/internal/entity"
	myError "github.com/Ararat25/subscription-aggregation-service/internal/error"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestReadSubscription - тест для функции ReadSubscription контроллера
func TestReadSubscription(t *testing.T) {
	mockService := new(MockAggregationService)

	handler := NewHandler(mockService)

	// Тестовый случай 1: Успешное чтение подписки
	{
		expectedSub := &entity.Subscription{
			Id:          1,
			ServiceName: "Test Service",
			Price:       100,
			UserId:      uuid.New(),
			StartDate:   time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC),
		}

		req := httptest.NewRequest("GET", "/subscription/{id}", nil)
		rw := httptest.NewRecorder()

		routeContext := chi.NewRouteContext()
		routeContext.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

		mockService.On("ReadSubscription", mock.Anything, int64(1)).Return(expectedSub, nil).Once()

		handler.ReadSubscription(rw, req)

		assert.Equal(t, http.StatusOK, rw.Code)
		var actualSub entity.SubscriptionRequest
		_ = json.Unmarshal(rw.Body.Bytes(), &actualSub)

		assert.Equal(t, expectedSub.Id, actualSub.Id)
		assert.Equal(t, expectedSub.ServiceName, actualSub.ServiceName)
		assert.Equal(t, expectedSub.Price, actualSub.Price)
		assert.Equal(t, expectedSub.UserId, actualSub.UserId)

		actualStartTime, err1 := time.Parse(entity.DateLayout, actualSub.StartDate)
		assert.NoError(t, err1)
		assert.WithinDuration(t, expectedSub.StartDate, actualStartTime, time.Second)
		mockService.AssertExpectations(t)
	}

	// Тестовый случай 2: Отсутствует параметр ID
	{
		req := httptest.NewRequest("GET", "/subscription/{id}", nil)
		rw := httptest.NewRecorder()

		handler.ReadSubscription(rw, req)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		var errResp ErrorResponse
		_ = json.Unmarshal(rw.Body.Bytes(), &errResp)
		assert.Contains(t, errResp.Error, "id parameter not set")
		mockService.AssertNotCalled(t, "ReadSubscription")
	}

	// Тестовый случай 3: Некорректный параметр ID (не число)
	{
		req := httptest.NewRequest("GET", "/subscription/{id}", nil)
		rw := httptest.NewRecorder()

		routeContext := chi.NewRouteContext()
		routeContext.URLParams.Add("id", "abc")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

		handler.ReadSubscription(rw, req)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		var errResp ErrorResponse
		_ = json.Unmarshal(rw.Body.Bytes(), &errResp)
		assert.Contains(t, errResp.Error, "invalid id parameter")
		mockService.AssertNotCalled(t, "ReadSubscription")
	}

	// Тестовый случай 4: Подписка не найдена (ошибка myError.ErrSubscriptionNotFound)
	{
		req := httptest.NewRequest("GET", "/subscription/{id}", nil)
		rw := httptest.NewRecorder()

		routeContext := chi.NewRouteContext()
		routeContext.URLParams.Add("id", "2")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

		mockService.On("ReadSubscription", mock.Anything, int64(2)).Return(&entity.Subscription{}, myError.ErrSubscriptionNotFound).Once()

		handler.ReadSubscription(rw, req)

		assert.Equal(t, http.StatusNotFound, rw.Code)
		var errResp ErrorResponse
		_ = json.Unmarshal(rw.Body.Bytes(), &errResp)
		assert.Contains(t, errResp.Error, myError.ErrSubscriptionNotFound.Error())
		mockService.AssertExpectations(t)
	}

	// Тестовый случай 5: Внутренняя ошибка сервиса
	{
		req := httptest.NewRequest("GET", "/subscription/{id}", nil)
		rw := httptest.NewRecorder()

		routeContext := chi.NewRouteContext()
		routeContext.URLParams.Add("id", "3")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

		mockService.On("ReadSubscription", mock.Anything, int64(3)).Return(&entity.Subscription{}, errors.New("some internal error")).Once()

		handler.ReadSubscription(rw, req)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		var errResp ErrorResponse
		_ = json.Unmarshal(rw.Body.Bytes(), &errResp)
		assert.Contains(t, errResp.Error, "some internal error")
		mockService.AssertExpectations(t)
	}
}
