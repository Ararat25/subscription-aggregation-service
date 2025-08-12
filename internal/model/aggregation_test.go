package model

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Ararat25/subscription-aggregation-service/internal/entity"
	myError "github.com/Ararat25/subscription-aggregation-service/internal/error"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepo это mock реализация интерфейса репозитория
type MockRepo struct {
	mock.Mock
}

// ConnectDB имитирует соединение с бд
func (m *MockRepo) ConnectDB(ctx context.Context, dbHost string, dbUser string, dbPassword string, dbName string, dbPort int) error {
	args := m.Called(ctx, dbHost, dbUser, dbPassword, dbName, dbPort)
	return args.Error(0)
}

// CreateSubscription имитирует создание подписки
func (m *MockRepo) CreateSubscription(ctx context.Context, s *entity.Subscription) (int64, error) {
	args := m.Called(ctx, s)
	return args.Get(0).(int64), args.Error(1)
}

// ReadSubscription имитирует получение подписки по id
func (m *MockRepo) ReadSubscription(ctx context.Context, id int64) (*entity.Subscription, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entity.Subscription), args.Error(1)
}

// UpdateSubscription имитирует обновление подписки
func (m *MockRepo) UpdateSubscription(ctx context.Context, s *entity.Subscription) error {
	args := m.Called(ctx, s)
	return args.Error(0)
}

// DeleteSubscription имитирует удаление подписки
func (m *MockRepo) DeleteSubscription(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ListSubscriptions имитирует вывод списка подписок
func (m *MockRepo) ListSubscriptions(ctx context.Context) ([]*entity.Subscription, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entity.Subscription), args.Error(1)
}

// TotalCost имитирует вывод суммарной стоимости подписок
func (m *MockRepo) TotalCost(ctx context.Context, from, to time.Time, userID *uuid.UUID, serviceName *string) (int, error) {
	args := m.Called(ctx, from, to, userID, serviceName)
	return args.Int(0), args.Error(1)
}

// Close имитирует закрытие соединенеия с бд
func (m *MockRepo) Close(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// TestNewAggregationService тестирует создание сервиса для агрегации
func TestNewAggregationService(t *testing.T) {
	mockRepo := new(MockRepo)
	service := NewAggregationService(mockRepo)
	assert.NotNil(t, service)
	assert.Equal(t, mockRepo, service.Storage)
}

// TestCreateSubscription тестирует создание подписки
func TestCreateSubscription(t *testing.T) {
	mockRepo := new(MockRepo)
	service := NewAggregationService(mockRepo)
	ctx := context.Background()

	// Тестовый пример 1: Успешное создание
	subReq := &entity.SubscriptionRequest{
		ServiceName: "Test Service",
		Price:       100,
		UserId:      uuid.New(),
		StartDate:   "01-2023",
	}
	expectedSub := &entity.Subscription{
		ServiceName: "Test Service",
		Price:       100,
		UserId:      subReq.UserId,
		StartDate:   time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     nil,
	}
	mockRepo.On("CreateSubscription", ctx, expectedSub).Return(int64(1), nil).Once()
	id, err := service.CreateSubscription(ctx, subReq)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)
	mockRepo.AssertExpectations(t)

	// Тестовый пример 2: Недопустимый аргумент (нулевой запрос)
	id, err = service.CreateSubscription(ctx, nil)
	assert.Error(t, err)
	assert.Equal(t, int64(0), id)
	assert.Contains(t, err.Error(), "invalid argument error")

	// Тестовый пример 3: Неверный формат даты
	subReqInvalidDate := &entity.SubscriptionRequest{
		ServiceName: "Test Service",
		Price:       100,
		UserId:      uuid.New(),
		StartDate:   "invalid-date",
	}
	id, err = service.CreateSubscription(ctx, subReqInvalidDate)
	assert.Error(t, err)
	assert.Equal(t, int64(0), id)
	assert.Contains(t, err.Error(), "invalid date format")

	// Тестовый пример 4: Дата окончания предшествует дате начала
	subReqInvalidEndDate := &entity.SubscriptionRequest{
		ServiceName: "Test Service",
		Price:       100,
		UserId:      uuid.New(),
		StartDate:   "01-2024",
		EndDate:     func() *string { s := "01-2023"; return &s }(),
	}
	id, err = service.CreateSubscription(ctx, subReqInvalidEndDate)
	assert.Error(t, err)
	assert.Equal(t, int64(0), id)
	assert.Contains(t, err.Error(), myError.ErrDateRange.Error())

	// Тестовый пример 5: Ошибка репозитория
	subReqRepoError := &entity.SubscriptionRequest{
		ServiceName: "Test Service",
		Price:       100,
		UserId:      uuid.New(),
		StartDate:   "01-2023",
	}
	expectedSubRepoError := &entity.Subscription{
		ServiceName: "Test Service",
		Price:       100,
		UserId:      subReqRepoError.UserId,
		StartDate:   time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     nil,
	}
	mockRepo.On("CreateSubscription", ctx, expectedSubRepoError).Return(int64(0), errors.New("db error")).Once()
	id, err = service.CreateSubscription(ctx, subReqRepoError)
	assert.Error(t, err)
	assert.Equal(t, int64(0), id)
	assert.Contains(t, err.Error(), "db error")
	mockRepo.AssertExpectations(t)
}

// TestReadSubscription тестирует чтение подписки
func TestReadSubscription(t *testing.T) {
	mockRepo := new(MockRepo)
	service := NewAggregationService(mockRepo)
	ctx := context.Background()

	// Тестовый пример 1: Успешное чтение
	expectedSub := &entity.Subscription{
		Id:          1,
		ServiceName: "Test Service",
		Price:       100,
		UserId:      uuid.New(),
		StartDate:   time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
	mockRepo.On("ReadSubscription", ctx, int64(1)).Return(expectedSub, nil).Once()
	sub, err := service.ReadSubscription(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedSub, sub)
	mockRepo.AssertExpectations(t)

	// Тестовый пример 2: Нет подписки
	mockRepo.On("ReadSubscription", ctx, int64(2)).Return(&entity.Subscription{}, errors.New("not found")).Once()
	sub, err = service.ReadSubscription(ctx, 2)
	assert.Error(t, err)
	assert.Nil(t, sub)
	assert.Contains(t, err.Error(), "not found")
	mockRepo.AssertExpectations(t)
}

// TestUpdateSubscription тестирует обновление подписки
func TestUpdateSubscription(t *testing.T) {
	mockRepo := new(MockRepo)
	service := NewAggregationService(mockRepo)
	ctx := context.Background()

	// Тестовый пример 1: Успешное обновление
	subReq := &entity.SubscriptionRequest{
		Id:          1,
		ServiceName: "Updated Service",
		Price:       200,
		UserId:      uuid.New(),
		StartDate:   "01-2023",
	}
	expectedSub := &entity.Subscription{
		Id:          1,
		ServiceName: "Updated Service",
		Price:       200,
		UserId:      subReq.UserId,
		StartDate:   time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     nil,
	}
	mockRepo.On("UpdateSubscription", ctx, expectedSub).Return(nil).Once()
	err := service.UpdateSubscription(ctx, subReq)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)

	// Тестовый пример 2: Недопустимый аргумент (нулевой запрос)
	err = service.UpdateSubscription(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid argument error")

	// Тестовый пример 3: Неверный формат даты
	subReqInvalidDate := &entity.SubscriptionRequest{
		Id:          1,
		ServiceName: "Updated Service",
		Price:       200,
		UserId:      uuid.New(),
		StartDate:   "invalid-date",
	}
	err = service.UpdateSubscription(ctx, subReqInvalidDate)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid date format")

	// Тестовый пример 4: Дата окончания предшествует дате начала
	subReqInvalidEndDate := &entity.SubscriptionRequest{
		Id:          1,
		ServiceName: "Updated Service",
		Price:       200,
		UserId:      uuid.New(),
		StartDate:   "01-2024",
		EndDate:     func() *string { s := "01-2023"; return &s }(),
	}
	err = service.UpdateSubscription(ctx, subReqInvalidEndDate)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), myError.ErrDateRange.Error())

	// Тестовый пример 5: Ошибка репозитория
	subReqRepoError := &entity.SubscriptionRequest{
		Id:          1,
		ServiceName: "Updated Service",
		Price:       200,
		UserId:      uuid.New(),
		StartDate:   "01-2023",
	}
	expectedSubRepoError := &entity.Subscription{
		Id:          1,
		ServiceName: "Updated Service",
		Price:       200,
		UserId:      subReqRepoError.UserId,
		StartDate:   time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     nil,
	}
	mockRepo.On("UpdateSubscription", ctx, expectedSubRepoError).Return(errors.New("db error")).Once()
	err = service.UpdateSubscription(ctx, subReqRepoError)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
	mockRepo.AssertExpectations(t)
}

// TestDeleteSubscription тестирует удаление подписки
func TestDeleteSubscription(t *testing.T) {
	mockRepo := new(MockRepo)
	service := NewAggregationService(mockRepo)
	ctx := context.Background()

	// Тестовый пример 1: Успешное удаление
	mockRepo.On("DeleteSubscription", ctx, int64(1)).Return(nil).Once()
	err := service.DeleteSubscription(ctx, 1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)

	// Тестовый пример 2: Ошибка репозитория
	mockRepo.On("DeleteSubscription", ctx, int64(2)).Return(errors.New("db error")).Once()
	err = service.DeleteSubscription(ctx, 2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
	mockRepo.AssertExpectations(t)
}

// TestListSubscriptions тестирует вывод всех подписок
func TestListSubscriptions(t *testing.T) {
	mockRepo := new(MockRepo)
	service := NewAggregationService(mockRepo)
	ctx := context.Background()

	// Тестовый пример 1: Успешный вывод
	expectedSubs := []*entity.Subscription{
		{
			Id:          1,
			ServiceName: "Service 1",
			Price:       100,
			UserId:      uuid.New(),
			StartDate:   time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Id:          2,
			ServiceName: "Service 2",
			Price:       200,
			UserId:      uuid.New(),
			StartDate:   time.Date(2023, time.February, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	mockRepo.On("ListSubscriptions", ctx).Return(expectedSubs, nil).Once()
	subs, err := service.ListSubscriptions(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expectedSubs, subs)
	mockRepo.AssertExpectations(t)

	// Тестовый пример 2: Ошибка репозитория
	mockRepo.On("ListSubscriptions", ctx).Return([]*entity.Subscription{}, errors.New("db error")).Once()
	subs, err = service.ListSubscriptions(ctx)
	assert.Error(t, err)
	assert.Nil(t, subs)
	assert.Contains(t, err.Error(), "db error")
	mockRepo.AssertExpectations(t)
}

// TestTotalCost тестирует вывод суммарной стоимости подписок
func TestTotalCost(t *testing.T) {
	mockRepo := new(MockRepo)
	service := NewAggregationService(mockRepo)
	ctx := context.Background()

	from := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2023, time.December, 31, 0, 0, 0, 0, time.UTC)
	userID := uuid.New()
	serviceName := "Test Service"

	// Тестовый пример 1: Успешный расчет общей стоимости
	mockRepo.On("TotalCost", ctx, resetDay(from), resetDay(to), &userID, &serviceName).Return(1000, nil).Once()
	cost, err := service.TotalCost(ctx, from, to, &userID, &serviceName)
	assert.NoError(t, err)
	assert.Equal(t, 1000, cost)
	mockRepo.AssertExpectations(t)

	// Тестовый пример 2: Недопустимый диапазон дат (до < от)
	cost, err = service.TotalCost(ctx, to, from, &userID, &serviceName)
	assert.Error(t, err)
	assert.Equal(t, 0, cost)
	assert.Contains(t, err.Error(), "to must be >= from")

	// Тестовый пример 3: Ошибка репозитория
	mockRepo.On("TotalCost", ctx, resetDay(from), resetDay(to), (*uuid.UUID)(nil), (*string)(nil)).Return(0, errors.New("db error")).Once()
	cost, err = service.TotalCost(ctx, from, to, nil, nil)
	assert.Error(t, err)
	assert.Equal(t, 0, cost)
	assert.Contains(t, err.Error(), "db error")
	mockRepo.AssertExpectations(t)
}
