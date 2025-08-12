package persistent_test

import (
	"context"
	"os"
	"testing"

	"github.com/andreyxaxa/order_svc/internal/entity"
	"github.com/andreyxaxa/order_svc/internal/repo/persistent"
	"github.com/andreyxaxa/order_svc/pkg/postgres"
	"github.com/stretchr/testify/assert"
)

func setupRepo(t *testing.T) *persistent.OrdersRepo {
	url := os.Getenv("TEST_PG_URL")

	if url == "" {
		t.Fatal("TEST_PG_URL variable required")
	}

	pg, err := postgres.New(url)
	if err != nil {
		t.Fatalf("failed to connect to postgres: %v", err)
	}

	return persistent.New(pg)
}

func TestStoreAndGet(t *testing.T) {
	ctx := context.Background()
	repo := setupRepo(t)

	orderUID := "b563feb7b2b84b6test"

	order := entity.Order{
		OrderUID:          orderUID,
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

	err := repo.Store(ctx, order)
	assert.NoError(t, err)

	got, err := repo.GetOrder(ctx, orderUID)
	assert.NoError(t, err)

	assert.Equal(t, order.OrderUID, got.OrderUID)
	assert.Equal(t, order.TrackNumber, got.TrackNumber)
	assert.Equal(t, order.Entry, got.Entry)
	assert.Equal(t, order.Locale, got.Locale)
	assert.Equal(t, order.InternalSignature, got.InternalSignature)
	assert.Equal(t, order.CustomerID, got.CustomerID)
	assert.Equal(t, order.DeliveryService, got.DeliveryService)
	assert.Equal(t, order.ShardKey, got.ShardKey)
	assert.Equal(t, order.SmID, got.SmID)
	assert.Equal(t, order.OofShard, got.OofShard)
	assert.Equal(t, order.Delivery, got.Delivery)
	assert.Equal(t, order.Items, order.Items)
	assert.Equal(t, order.Payment, got.Payment)
}
