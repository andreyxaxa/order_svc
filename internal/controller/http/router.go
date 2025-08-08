package http

import (
	v1 "github.com/andreyxaxa/order_svc/internal/controller/http/v1"
	"github.com/andreyxaxa/order_svc/internal/usecase"
	"github.com/andreyxaxa/order_svc/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

func NewRouter(app *fiber.App, config, o usecase.Orders, l logger.Interface) {
	apiV1Group := app.Group("/v1")
	{
		v1.NewOrderRoutes(apiV1Group, o, l)
	}
}
