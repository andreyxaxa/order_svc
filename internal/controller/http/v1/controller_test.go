package v1_test

import (
	"context"
	"encoding/json"
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

func TestController_Success(t *testing.T) {
	app := fiber.New()
	mockUC := new(mockOrdersUseCase)

	expected := entity.Order{
		OrderUID:          "b563feb7b2b84b6test",
		TrackNumber:       "WBILMTESTTRACK",
		Entry:             "WBIL",
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test",
		DeliveryService:   "meest",
		ShardKey:          "9",
		SmID:              99,
		OofShard:          "1",
		Delivery: entity.Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: entity.Payment{
			Transaction:  "b563feb7b2b84b6test",
			RequestID:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDT:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items: []entity.Item{
			{
				ChrtID:      9934930,
				TrackNumber: "WBILMTESTTRACK",
				Price:       453,
				RID:         "ab4219087a764ae0btest",
				Name:        "Mascaras",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NmID:        2389212,
				Brand:       "Vivienne Sabo",
				Status:      202,
			},
		},
	}
	mockUC.On("Order", mock.Anything, "b563feb7b2b84b6test").Return(expected, nil)

	v1.NewOrderRoutes(app.Group("/v1"), mockUC, nil)

	req := httptest.NewRequest(http.MethodGet, "/v1/order/info?order_uid=b563feb7b2b84b6test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	actual := entity.Order{}
	err = json.NewDecoder(resp.Body).Decode(&actual)
	assert.NoError(t, err)

	assert.Equal(t, actual, expected)

	mockUC.AssertExpectations(t)
}
