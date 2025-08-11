package controller

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	myError "github.com/Ararat25/subscription-aggregation-service/internal/error"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestDeleteSubscription - тест для функции DeleteSubscription контроллера
func TestDeleteSubscription(t *testing.T) {
	mockService := new(MockAggregationService)

	handler := NewHandler(mockService)

	// Тестовый случай 1: Успешное удаление подписки
	{
		req := httptest.NewRequest("DELETE", "/subscription/delete/{id}", nil)
		rw := httptest.NewRecorder()

		routeContext := chi.NewRouteContext()
		routeContext.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

		mockService.On("DeleteSubscription", mock.Anything, int64(1)).Return(nil).Once()

		handler.DeleteSubscription(rw, req)

		assert.Equal(t, http.StatusOK, rw.Code)
		var resp StatusResponse
		_ = json.Unmarshal(rw.Body.Bytes(), &resp)
		assert.Equal(t, "success", resp.Status)
		mockService.AssertExpectations(t)
	}

	// Тестовый случай 2: Отсутствует параметр ID
	{
		req := httptest.NewRequest("DELETE", "/subscription/delete/{id}", nil)
		rw := httptest.NewRecorder()

		handler.DeleteSubscription(rw, req)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		var errResp ErrorResponse
		_ = json.Unmarshal(rw.Body.Bytes(), &errResp)
		assert.Contains(t, errResp.Error, "id parameter not set")
		mockService.AssertNotCalled(t, "DeleteSubscription")
	}

	// Тестовый случай 3: Некорректный параметр ID (не число)
	{
		req := httptest.NewRequest("DELETE", "/subscription/delete/{id}", nil)
		rw := httptest.NewRecorder()

		routeContext := chi.NewRouteContext()
		routeContext.URLParams.Add("id", "abc")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

		handler.DeleteSubscription(rw, req)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		var errResp ErrorResponse
		_ = json.Unmarshal(rw.Body.Bytes(), &errResp)
		assert.Contains(t, errResp.Error, "invalid id parameter")
		mockService.AssertNotCalled(t, "DeleteSubscription")
	}

	// Тестовый случай 4: Подписка не найдена (ошибка myError.ErrSubscriptionNotFound)
	{
		req := httptest.NewRequest("DELETE", "/subscription/delete/{id}", nil)
		rw := httptest.NewRecorder()

		routeContext := chi.NewRouteContext()
		routeContext.URLParams.Add("id", "2")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

		mockService.On("DeleteSubscription", mock.Anything, int64(2)).Return(myError.ErrSubscriptionNotFound).Once()

		handler.DeleteSubscription(rw, req)

		assert.Equal(t, http.StatusNotFound, rw.Code)
		var errResp ErrorResponse
		_ = json.Unmarshal(rw.Body.Bytes(), &errResp)
		assert.Contains(t, errResp.Error, myError.ErrSubscriptionNotFound.Error())
		mockService.AssertExpectations(t)
	}

	// Тестовый случай 5: Внутренняя ошибка сервиса
	{
		req := httptest.NewRequest("DELETE", "/subscription/delete/{id}", nil)
		rw := httptest.NewRecorder()

		routeContext := chi.NewRouteContext()
		routeContext.URLParams.Add("id", "3")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

		mockService.On("DeleteSubscription", mock.Anything, int64(3)).Return(errors.New("some internal error")).Once()

		handler.DeleteSubscription(rw, req)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		var errResp ErrorResponse
		_ = json.Unmarshal(rw.Body.Bytes(), &errResp)
		assert.Contains(t, errResp.Error, "some internal error")
		mockService.AssertExpectations(t)
	}
}
