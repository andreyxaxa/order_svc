package orders_test

import (
	"context"
	"errors"
	"testing"

	"github.com/andreyxaxa/order_svc/internal/entity"
	"github.com/andreyxaxa/order_svc/internal/usecase/orders"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockOrdersRepo struct {
	mock.Mock
}

func (m *mockOrdersRepo) GetOrder(ctx context.Context, orderUID string) (entity.Order, error) {
	args := m.Called(ctx, orderUID)
	return args.Get(0).(entity.Order), args.Error(1)
}

func (m *mockOrdersRepo) Store(ctx context.Context, order entity.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func TestOrderUseCase_Success(t *testing.T) {
	mockRepo := new(mockOrdersRepo)
	expected := entity.Order{OrderUID: "123test", TrackNumber: "WBTEST"}

	mockRepo.On("GetOrder", mock.Anything, "123test").Return(expected, nil)

	uc := orders.New(mockRepo)
	got, err := uc.Order(context.Background(), "123test")

	assert.NoError(t, err)
	assert.Equal(t, expected, got)
	mockRepo.AssertExpectations(t)
}

func TestOrderUseCase_Error(t *testing.T) {
	mockRepo := new(mockOrdersRepo)

	mockRepo.On("GetOrder", mock.Anything, "123test").Return(entity.Order{}, errors.New("storage problems"))

	uc := orders.New(mockRepo)
	_, err := uc.Order(context.Background(), "123test")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "storage problems")
	mockRepo.AssertExpectations(t)
}
