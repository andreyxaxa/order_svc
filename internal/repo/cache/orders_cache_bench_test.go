package cache_test

import (
	"context"
	"os"
	"testing"

	"github.com/andreyxaxa/order_svc/internal/entity"
	"github.com/andreyxaxa/order_svc/internal/repo/cache"
	"github.com/andreyxaxa/order_svc/internal/repo/persistent"
	"github.com/andreyxaxa/order_svc/pkg/postgres"
)

func getExampleOrder() entity.Order {
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

	return order
}

func getRepos(tb testing.TB) (*cache.CachedOrdersRepo, *persistent.OrdersRepo) {
	url := os.Getenv("TEST_PG_URL")

	if url == "" {
		tb.Fatal("TEST_PG_URL variable required")
	}

	// postgres
	pg, err := postgres.New(url)
	if err != nil {
		tb.Fatalf("failed to connect to postgres: %v", err)
	}

	// postgres repo
	persistentRepo := persistent.New(pg)

	// cached repo
	cachedRepo := cache.New(persistentRepo, 100, 5)

	cachedRepo.Store(context.Background(), getExampleOrder())

	return cachedRepo, persistentRepo
}

func BenchmarkGet(b *testing.B) {
	ctx := context.Background()
	cachedRepo, persistentRepo := getRepos(b)

	orderUID := "b563feb7b2b84b6test"

	// без кеша
	b.Run("PersistentRepo", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := persistentRepo.GetOrder(ctx, orderUID)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	// с кешем
	b.Run("CachedPersistentRepo", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := cachedRepo.GetOrder(ctx, orderUID)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
