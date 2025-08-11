package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Ararat25/subscription-aggregation-service/internal/entity"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAggregationService - мок для интерфейса AggregationService
type MockAggregationService struct {
	mock.Mock
}

// CreateSubscription - мок метод для создания подписки
func (m *MockAggregationService) CreateSubscription(ctx context.Context, s *entity.SubscriptionRequest) (int64, error) {
	args := m.Called(ctx, s)
	return args.Get(0).(int64), args.Error(1)
}

// ReadSubscription - мок метод для чтения подписки
func (m *MockAggregationService) ReadSubscription(ctx context.Context, id int64) (*entity.Subscription, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entity.Subscription), args.Error(1)
}

// UpdateSubscription - мок метод для обновления подписки
func (m *MockAggregationService) UpdateSubscription(ctx context.Context, s *entity.SubscriptionRequest) error {
	args := m.Called(ctx, s)
	return args.Error(0)
}

// DeleteSubscription - мок метод для удаления подписки
func (m *MockAggregationService) DeleteSubscription(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ListSubscriptions - мок метод для получения списка подписок
func (m *MockAggregationService) ListSubscriptions(ctx context.Context) ([]*entity.Subscription, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entity.Subscription), args.Error(1)
}

// TotalCost - мок метод для подсчета общей стоимости подписок
func (m *MockAggregationService) TotalCost(ctx context.Context, from, to time.Time, userID *uuid.UUID, serviceName *string) (int, error) {
	args := m.Called(ctx, from, to, userID, serviceName)
	return args.Int(0), args.Error(1)
}

// TestCreateSubscription - тест для CreateSubscription контроллера
func TestCreateSubscription(t *testing.T) {
	validate = validator.New()

	mockService := new(MockAggregationService)

	handler := NewHandler(mockService)

	// Тестовый случай 1: Успешное создание подписки
	{
		subReq := entity.SubscriptionRequest{
			ServiceName: "Test Service",
			Price:       100,
			UserId:      uuid.New(),
			StartDate:   "01-2023",
		}
		jsonBody, _ := json.Marshal(subReq)
		req, _ := http.NewRequest("POST", "/subscription", bytes.NewBuffer(jsonBody))
		rw := httptest.NewRecorder()

		mockService.On("CreateSubscription", mock.Anything, mock.AnythingOfType("*entity.SubscriptionRequest")).Return(int64(1), nil).Once()

		handler.CreateSubscription(rw, req)

		assert.Equal(t, http.StatusOK, rw.Code)
		var resp CreateControllerResponse
		_ = json.Unmarshal(rw.Body.Bytes(), &resp)
		assert.Equal(t, int64(1), resp.Id)
		mockService.AssertExpectations(t)
	}

	// Тестовый случай 2: Некорректный JSON в теле запроса
	{
		req, _ := http.NewRequest("POST", "/subscription", bytes.NewBuffer([]byte("invalid json")))
		rw := httptest.NewRecorder()

		handler.CreateSubscription(rw, req)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		var errResp ErrorResponse
		_ = json.Unmarshal(rw.Body.Bytes(), &errResp)
		assert.Contains(t, errResp.Error, "invalid character")
		mockService.AssertNotCalled(t, "CreateSubscription")
	}

	// Тестовый случай 3: Ошибка валидации данных подписки
	{
		subReq := entity.SubscriptionRequest{
			Price:     100,
			UserId:    uuid.New(),
			StartDate: "01-2023",
		}
		jsonBody, _ := json.Marshal(subReq)
		req, _ := http.NewRequest("POST", "/subscription", bytes.NewBuffer(jsonBody))
		rw := httptest.NewRecorder()

		handler.CreateSubscription(rw, req)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		var errResp ErrorResponse
		_ = json.Unmarshal(rw.Body.Bytes(), &errResp)
		assert.Contains(t, errResp.Error, "ServiceName")
		mockService.AssertNotCalled(t, "CreateSubscription")
	}

	// Тестовый случай 4: Ошибка сервиса агрегации при создании подписки
	{
		subReq := entity.SubscriptionRequest{
			ServiceName: "Test Service",
			Price:       100,
			UserId:      uuid.New(),
			StartDate:   "01-2023",
		}
		jsonBody, _ := json.Marshal(subReq)
		req, _ := http.NewRequest("POST", "/subscription", bytes.NewBuffer(jsonBody))
		rw := httptest.NewRecorder()

		mockService.On("CreateSubscription", mock.Anything, mock.AnythingOfType("*entity.SubscriptionRequest")).Return(int64(0), errors.New("internal service error")).Once()

		handler.CreateSubscription(rw, req)

		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		var errResp ErrorResponse
		_ = json.Unmarshal(rw.Body.Bytes(), &errResp)
		assert.Contains(t, errResp.Error, "internal service error")
		mockService.AssertExpectations(t)
	}
}
