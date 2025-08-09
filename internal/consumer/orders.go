package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/andreyxaxa/order_svc/internal/entity"
	"github.com/andreyxaxa/order_svc/internal/repo"
	"github.com/andreyxaxa/order_svc/pkg/kafka"
	"github.com/andreyxaxa/order_svc/pkg/logger"
)

type OrdersConsumer struct {
	k      *kafka.Kafka
	r      repo.OrdersRepo
	l      logger.Interface
	cancel context.CancelFunc
}

func New(k *kafka.Kafka, r repo.OrdersRepo, l logger.Interface) *OrdersConsumer {
	return &OrdersConsumer{
		k: k,
		r: r,
		l: l,
	}
}

func (c *OrdersConsumer) Start(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	c.cancel = cancel

	defer func() {
		if err := c.k.Close(); err != nil {
			c.l.Error(err, "Consumer - Start - c.k.Close")
		}
	}()

	for {
		select {
		case <-ctx.Done():
			c.l.Info("Consumer stopped by context cancel")
			return
		default:
			msg, err := c.k.ReadMessage(ctx)
			if err != nil {
				c.l.Error(err, "Consumer - Start - c.k.ReadMessage")
				continue
			}

			var order entity.Order
			if err := json.Unmarshal(msg.Value, &order); err != nil {
				c.l.Error(err, "Consumer - Start - json.Unmarshal")
				continue
			}

			if err := c.r.Store(ctx, order); err != nil {
				c.l.Error(err, "Consumer - Start - c.r.Store")
				continue
			}

			c.l.Info(fmt.Sprintf("Order %s stored", order.OrderUID))
		}
	}
}

func (c *OrdersConsumer) Stop() {
	if c.cancel != nil {
		c.cancel()
	}
}
