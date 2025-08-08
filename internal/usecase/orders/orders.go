package orders

import (
	"context"
	"fmt"

	"github.com/andreyxaxa/order_svc/internal/entity"
	"github.com/andreyxaxa/order_svc/internal/repo"
)

type UseCase struct {
	repo repo.OrdersRepo
}

func New(r repo.OrdersRepo) *UseCase {
	return &UseCase{
		repo: r,
	}
}

func (uc *UseCase) Order(ctx context.Context, orderUID string) (entity.Order, error) {
	order, err := uc.repo.GetOrder(ctx, orderUID)
	if err != nil {
		return entity.Order{}, fmt.Errorf("OrderUseCase - Order - uc.repo.GetOrder: %w", err)
	}

	return order, nil
}
