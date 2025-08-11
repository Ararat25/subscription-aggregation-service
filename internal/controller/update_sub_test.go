package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ararat25/subscription-aggregation-service/internal/entity"
	myError "github.com/Ararat25/subscription-aggregation-service/internal/error"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestUpdateSubscription - тест для функции UpdateSubscription контроллера
func TestUpdateSubscription(t *testing.T) {
	validate = validator.New()

	mockService := new(MockAggregationService)

	handler := NewHandler(mockService)

	// Тестовый случай 1: Успешное обновление подписки
	{
		subReq := entity.SubscriptionRequest{
			Id:          1,
			ServiceName: "Updated Service",
			Price:       200,
			UserId:      uuid.New(),
			StartDate:   "01-2023",
		}
		jsonBody, _ := json.Marshal(subReq)
		req, _ := http.NewRequest("PUT", "/subscription/update", bytes.NewBuffer(jsonBody))
		rw := httptest.NewRecorder()

		mockService.On("UpdateSubscription", mock.Anything, mock.AnythingOfType("*entity.SubscriptionRequest")).Return(nil).Once()

		handler.UpdateSubscription(rw, req)

		assert.Equal(t, http.StatusOK, rw.Code)
		var resp StatusResponse
		_ = json.Unmarshal(rw.Body.Bytes(), &resp)
		assert.Equal(t, "success", resp.Status)
		mockService.AssertExpectations(t)
	}

	// Тестовый случай 2: Некорректный JSON в теле запроса
	{
		req, _ := http.NewRequest("PUT", "/subscription/update", bytes.NewBuffer([]byte("invalid json")))
		rw := httptest.NewRecorder()

		handler.UpdateSubscription(rw, req)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		var errResp ErrorResponse
		_ = json.Unmarshal(rw.Body.Bytes(), &errResp)
		assert.Contains(t, errResp.Error, "invalid character")
		mockService.AssertNotCalled(t, "UpdateSubscription")
	}

	// Тестовый случай 3: Ошибка валидации данных подписки
	{
		subReq := entity.SubscriptionRequest{
			Id:        1,
			Price:     200,
			UserId:    uuid.New(),
			StartDate: "01-2023",
		}
		jsonBody, _ := json.Marshal(subReq)
		req, _ := http.NewRequest("PUT", "/subscription/update", bytes.NewBuffer(jsonBody))
		rw := httptest.NewRecorder()

		handler.UpdateSubscription(rw, req)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		var errResp ErrorResponse
		_ = json.Unmarshal(rw.Body.Bytes(), &errResp)
		assert.Contains(t, errResp.Error, "ServiceName")
		mockService.AssertNotCalled(t, "UpdateSubscription")
	}

	// Тестовый случай 4: Подписка не найдена (ошибка myError.ErrSubscriptionNotFound)
	{
		subReq := entity.SubscriptionRequest{
			Id:          999,
			ServiceName: "Non Existent Service",
			Price:       200,
			UserId:      uuid.New(),
			StartDate:   "01-2023",
		}
		jsonBody, _ := json.Marshal(subReq)
		req, _ := http.NewRequest("PUT", "/subscription/update", bytes.NewBuffer(jsonBody))
		rw := httptest.NewRecorder()

		mockService.On("UpdateSubscription", mock.Anything, mock.AnythingOfType("*entity.SubscriptionRequest")).Return(myError.ErrSubscriptionNotFound).Once()

		handler.UpdateSubscription(rw, req)

		assert.Equal(t, http.StatusNotFound, rw.Code)
		var errResp ErrorResponse
		_ = json.Unmarshal(rw.Body.Bytes(), &errResp)
		assert.Contains(t, errResp.Error, myError.ErrSubscriptionNotFound.Error())
		mockService.AssertExpectations(t)
	}

	// Тестовый случай 5: Внутренняя ошибка сервиса
	{
		subReq := entity.SubscriptionRequest{
			Id:          1,
			ServiceName: "Test Service",
			Price:       100,
			UserId:      uuid.New(),
			StartDate:   "01-2023",
		}
		jsonBody, _ := json.Marshal(subReq)
		req, _ := http.NewRequest("PUT", "/subscription/update", bytes.NewBuffer(jsonBody))
		rw := httptest.NewRecorder()

		mockService.On("UpdateSubscription", mock.Anything, mock.AnythingOfType("*entity.SubscriptionRequest")).Return(errors.New("some internal error")).Once()

		handler.UpdateSubscription(rw, req)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		var errResp ErrorResponse
		_ = json.Unmarshal(rw.Body.Bytes(), &errResp)
		assert.Contains(t, errResp.Error, "some internal error")
		mockService.AssertExpectations(t)
	}
}
