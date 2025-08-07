package repo

import (
	"context"

	"github.com/andreyxaxa/order_svc/internal/entity"
)

type OrdersRepo interface {
	Store(context.Context, entity.Order) error
	GetOrder(context.Context, string) (entity.Order, error)
}
