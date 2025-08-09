package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/andreyxaxa/order_svc/config"

	"github.com/andreyxaxa/order_svc/internal/consumer"
	"github.com/andreyxaxa/order_svc/internal/controller/http"
	"github.com/andreyxaxa/order_svc/internal/repo/persistent"
	"github.com/andreyxaxa/order_svc/internal/usecase/orders"
	"github.com/andreyxaxa/order_svc/pkg/httpserver"
	"github.com/andreyxaxa/order_svc/pkg/kafka"
	"github.com/andreyxaxa/order_svc/pkg/logger"
	"github.com/andreyxaxa/order_svc/pkg/postgres"
)

func Run(cfg *config.Config) {
	// Logger
	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()
	repo := persistent.New(pg)

	// Use-Case
	ordersUseCase := orders.New(repo)

	// Kafka Consumer
	ordersConsumer := consumer.New(kafka.New(kafka.Topic(cfg.Kafka.Topic)), repo, l)

	// HTTP Server
	httpServer := httpserver.New(httpserver.Port(cfg.HTTP.Port))
	http.NewRouter(httpServer.App, ordersUseCase, l)

	// Consumer Start
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go ordersConsumer.Start(ctx)

	// Server Start
	httpServer.Start()

	// Waiting Signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
	case err := <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
		cancel()
	}

	ordersConsumer.Stop()

	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
