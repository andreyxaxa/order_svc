package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	_defaultBrokerAddress = "broker:19092"
	_defaultTopic         = "orders"
	_defaultGroupID       = "order-consumer"
	_defaultMinBytes      = 10e3
	_defaultMaxBytes      = 10e6
	_defaultMaxAttempts   = 10
	_defaultReadTimeout   = 2 * time.Minute
)

type Kafka struct {
	brokers      []string
	topic        string
	groupID      string
	minBytes     int
	maxBytes     int
	readTimeout  time.Duration
	connAttempts int

	Reader *kafka.Reader
}

func New(opts ...Option) *Kafka {
	k := &Kafka{
		brokers:      []string{_defaultBrokerAddress},
		topic:        _defaultTopic,
		groupID:      _defaultGroupID,
		minBytes:     _defaultMinBytes,
		maxBytes:     _defaultMaxBytes,
		readTimeout:  _defaultReadTimeout,
		connAttempts: _defaultMaxAttempts,
	}

	// Custom options
	for _, opt := range opts {
		opt(k)
	}

	k.Reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:     k.brokers,
		GroupID:     k.groupID,
		Topic:       k.topic,
		MinBytes:    k.minBytes,
		MaxBytes:    k.maxBytes,
		MaxAttempts: k.connAttempts,
	})

	return k
}

func (k *Kafka) ReadMessage(ctx context.Context) (kafka.Message, error) {
	ctx, cancel := context.WithTimeout(ctx, k.readTimeout)
	defer cancel()

	msg, err := k.Reader.ReadMessage(ctx)
	if err != nil {
		return kafka.Message{}, fmt.Errorf("kafka - ReadMessage - k.Reader.ReadMessage: %w", err)
	}

	return msg, err
}

func (k *Kafka) Close() error {
	if k.Reader != nil {
		return k.Reader.Close()
	}

	return nil
}
