package v1

import (
	"github.com/andreyxaxa/order_svc/internal/usecase"
	"github.com/andreyxaxa/order_svc/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

func NewOrderRoutes(apiV1Group fiber.Router, o usecase.Orders, l logger.Interface) {
	r := &V1{
		o: o,
		l: l,
	}

	orderGroup := apiV1Group.Group("/order")

	{
		orderGroup.Get("/info", r.orderJSON)
		orderGroup.Get("/info/html", r.orderHTML)
	}
}
