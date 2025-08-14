package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/andreyxaxa/order_svc/internal/entity"
	"github.com/andreyxaxa/order_svc/internal/repo/cache/lru"
	"github.com/andreyxaxa/order_svc/internal/repo/persistent"
	errs "github.com/andreyxaxa/order_svc/pkg/errors"
)

type CachedOrdersRepo struct {
	db    *persistent.OrdersRepo
	cache *lru.LRUCache
}

func New(dbRepo *persistent.OrdersRepo, capacity int, ttlMinutes int) *CachedOrdersRepo {
	return &CachedOrdersRepo{
		db:    dbRepo,
		cache: lru.New(capacity, time.Minute*time.Duration(ttlMinutes)),
	}
}

func (r *CachedOrdersRepo) GetOrder(ctx context.Context, orderUID string) (entity.Order, error) {
	if order, ok := r.cache.Get(orderUID); ok {
		return order, nil
	}

	order, err := r.db.GetOrder(ctx, orderUID)
	if err != nil {
		if errors.Is(err, errs.ErrNoRows) {
			return entity.Order{}, errs.ErrNoRows
		}
		return entity.Order{}, fmt.Errorf("CachedOrdersRepo - GetOrder - r.db.GetOrder: %w", err)
	}

	r.cache.Set(orderUID, order)

	return order, nil
}

func (r *CachedOrdersRepo) Store(ctx context.Context, order entity.Order) error {
	if err := r.db.Store(ctx, order); err != nil {
		return fmt.Errorf("CachedOrderRepo - Store - r.db.Store: %w", err)
	}

	r.cache.Set(order.OrderUID, order)

	return nil
}

func (r *CachedOrdersRepo) PreloadCache(ctx context.Context, limit int) error {
	orders, err := r.db.ListRecentOrders(ctx, limit)
	if err != nil {
		return fmt.Errorf("CachedOrdersRepo - PreloadCache - r.db.ListRecentOrders")
	}

	for _, order := range orders {
		r.cache.Set(order.OrderUID, order)
	}

	return nil
}
