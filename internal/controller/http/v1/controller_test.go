package v1_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	v1 "github.com/andreyxaxa/order_svc/internal/controller/http/v1"
	"github.com/andreyxaxa/order_svc/internal/entity"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockOrdersUseCase struct {
	mock.Mock
}

func (m *mockOrdersUseCase) Order(ctx context.Context, orderUID string) (entity.Order, error) {
	args := m.Called(ctx, orderUID)
	return args.Get(0).(entity.Order), args.Error(1)
}

// TODO: доделать тест(более обширно)
func TestController_Success(t *testing.T) {
	app := fiber.New()
	mockUC := new(mockOrdersUseCase)

	expected := entity.Order{OrderUID: "123test", TrackNumber: "WBTEST"}
	mockUC.On("Order", mock.Anything, "123test").Return(expected, nil)

	v1.NewOrderRoutes(app.Group("/v1"), mockUC, nil)

	req := httptest.NewRequest(http.MethodGet, "/v1/order/info?order_uid=123test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	bodyStr := string(bodyBytes)

	assert.Contains(t, bodyStr, expected.OrderUID)
	assert.Contains(t, bodyStr, expected.TrackNumber)

	mockUC.AssertExpectations(t)
}
