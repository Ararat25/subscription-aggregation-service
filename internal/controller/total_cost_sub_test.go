package controller

import (
	"encoding/json"
	"errors"
	"github.com/Ararat25/subscription-aggregation-service/internal/entity"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestTotalCost - тест для функции TotalCost контроллера
func TestTotalCost(t *testing.T) {
	mockService := new(MockAggregationService)

	handler := NewHandler(mockService)

	// Тестовый случай 1: Успешное вычисление общей стоимости со всеми параметрами
	{
		from := "01-2023"
		to := "12-2023"
		userID := uuid.New()
		serviceName := "Test Service"

		params := url.Values{}
		params.Add("from", from)
		params.Add("to", to)
		params.Add("id", userID.String())
		params.Add("service_name", serviceName)

		req := httptest.NewRequest("GET", "/subscriptions/cost?"+params.Encode(), nil)
		rw := httptest.NewRecorder()

		parsedFrom, _ := time.Parse(entity.DateLayout, from)
		parsedTo, _ := time.Parse(entity.DateLayout, to)

		mockService.On("TotalCost", mock.Anything, parsedFrom, parsedTo, &userID, &serviceName).Return(1000, nil).Once()

		handler.TotalCost(rw, req)

		assert.Equal(t, http.StatusOK, rw.Code)
		var resp TotalCostControllerResponse
		_ = json.Unmarshal(rw.Body.Bytes(), &resp)
		assert.Equal(t, 1000, resp.TotalCost)
		mockService.AssertExpectations(t)
	}

	// Тестовый случай 2: Успешное вычисление общей стоимости только с датами from и to
	{
		from := "01-2023"
		to := "12-2023"

		params := url.Values{}
		params.Add("from", from)
		params.Add("to", to)

		req := httptest.NewRequest("GET", "/subscriptions/cost?"+params.Encode(), nil)
		rw := httptest.NewRecorder()

		parsedFrom, _ := time.Parse(entity.DateLayout, from)
		parsedTo, _ := time.Parse(entity.DateLayout, to)

		mockService.On("TotalCost", mock.Anything, parsedFrom, parsedTo, (*uuid.UUID)(nil), (*string)(nil)).Return(500, nil).Once()

		handler.TotalCost(rw, req)

		assert.Equal(t, http.StatusOK, rw.Code)
		var resp TotalCostControllerResponse
		_ = json.Unmarshal(rw.Body.Bytes(), &resp)
		assert.Equal(t, 500, resp.TotalCost)
		mockService.AssertExpectations(t)
	}

	// Тестовый случай 3: Отсутствует параметр from или to
	{
		params := url.Values{}
		params.Add("from", "01-2023")

		req := httptest.NewRequest("GET", "/subscriptions/cost?"+params.Encode(), nil)
		rw := httptest.NewRecorder()

		handler.TotalCost(rw, req)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		var errResp ErrorResponse
		_ = json.Unmarshal(rw.Body.Bytes(), &errResp)
		assert.Contains(t, errResp.Error, "missing from or to parameter")
		mockService.AssertNotCalled(t, "TotalCost")
	}

	// Тестовый случай 4: Некорректный формат даты from
	{
		params := url.Values{}
		params.Add("from", "invalid-date")
		params.Add("to", "12-2023")

		req := httptest.NewRequest("GET", "/subscriptions/cost?"+params.Encode(), nil)
		rw := httptest.NewRecorder()

		handler.TotalCost(rw, req)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		var errResp ErrorResponse
		_ = json.Unmarshal(rw.Body.Bytes(), &errResp)
		assert.Contains(t, errResp.Error, "invalid from parameter")
		mockService.AssertNotCalled(t, "TotalCost")
	}

	// Тестовый случай 5: Некорректный формат даты to
	{
		params := url.Values{}
		params.Add("from", "01-2023")
		params.Add("to", "invalid-date")

		req := httptest.NewRequest("GET", "/subscriptions/cost?"+params.Encode(), nil)
		rw := httptest.NewRecorder()

		handler.TotalCost(rw, req)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		var errResp ErrorResponse
		_ = json.Unmarshal(rw.Body.Bytes(), &errResp)
		assert.Contains(t, errResp.Error, "invalid to parameter")
		mockService.AssertNotCalled(t, "TotalCost")
	}

	// Тестовый случай 6: Некорректный формат ID пользователя
	{
		params := url.Values{}
		params.Add("from", "01-2023")
		params.Add("to", "12-2023")
		params.Add("id", "invalid-uuid")

		req := httptest.NewRequest("GET", "/subscriptions/cost?"+params.Encode(), nil)
		rw := httptest.NewRecorder()

		handler.TotalCost(rw, req)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		var errResp ErrorResponse
		_ = json.Unmarshal(rw.Body.Bytes(), &errResp)
		assert.Contains(t, errResp.Error, "invalid id")
		mockService.AssertNotCalled(t, "TotalCost")
	}

	// Тестовый случай 7: Ошибка сервиса агрегации
	{
		from := "01-2023"
		to := "12-2023"

		params := url.Values{}
		params.Add("from", from)
		params.Add("to", to)

		req := httptest.NewRequest("GET", "/subscriptions/cost?"+params.Encode(), nil)
		rw := httptest.NewRecorder()

		parsedFrom, _ := time.Parse(entity.DateLayout, from)
		parsedTo, _ := time.Parse(entity.DateLayout, to)

		mockService.On("TotalCost", mock.Anything, parsedFrom, parsedTo, (*uuid.UUID)(nil), (*string)(nil)).Return(0, errors.New("internal service error")).Once()

		handler.TotalCost(rw, req)

		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		var errResp ErrorResponse
		_ = json.Unmarshal(rw.Body.Bytes(), &errResp)
		assert.Contains(t, errResp.Error, "internal service error")
		mockService.AssertExpectations(t)
	}
}
