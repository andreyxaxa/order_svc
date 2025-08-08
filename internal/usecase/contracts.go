package usecase

import (
	"context"

	"github.com/andreyxaxa/order_svc/internal/entity"
)

type Orders interface {
	Order(context.Context, string) (entity.Order, error)
}
