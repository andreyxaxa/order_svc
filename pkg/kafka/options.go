package kafka

import "time"

type Option func(k *Kafka)

func Brokers(brokers []string) Option {
	return func(k *Kafka) {
		k.brokers = brokers
	}
}

func Topic(topic string) Option {
	return func(k *Kafka) {
		k.topic = topic
	}
}

func GroupID(groupID string) Option {
	return func(k *Kafka) {
		k.groupID = groupID
	}
}

func MinBytes(minBytes int) Option {
	return func(k *Kafka) {
		k.minBytes = minBytes
	}
}

func MaxBytes(maxBytes int) Option {
	return func(k *Kafka) {
		k.maxBytes = maxBytes
	}
}

func ConnAttempts(attempts int) Option {
	return func(k *Kafka) {
		k.connAttempts = attempts
	}
}

func ReadTimeout(timeout time.Duration) Option {
	return func(k *Kafka) {
		k.readTimeout = timeout
	}
}
