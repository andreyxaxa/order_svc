package http

import (
	"github.com/andreyxaxa/order_svc/config"
	_ "github.com/andreyxaxa/order_svc/docs" // Swagger docs
	v1 "github.com/andreyxaxa/order_svc/internal/controller/http/v1"
	"github.com/andreyxaxa/order_svc/internal/usecase"
	"github.com/andreyxaxa/order_svc/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// @title Order-Check service
// @version 1.0
// @host localhost:8080
// @BasePath /v1
func NewRouter(app *fiber.App, cfg *config.Config, o usecase.Orders, l logger.Interface) {
	// Swagger
	if cfg.Swagger.Enabled {
		app.Get("/swagger/*", swagger.HandlerDefault)
	}

	// Routers
	apiV1Group := app.Group("/v1")
	{
		v1.NewOrderRoutes(apiV1Group, o, l)
	}
}
