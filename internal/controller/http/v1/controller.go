package v1

import (
	"github.com/andreyxaxa/order_svc/internal/usecase"
	"github.com/andreyxaxa/order_svc/pkg/logger"
)

type V1 struct {
	o usecase.Orders
	l logger.Interface
}
