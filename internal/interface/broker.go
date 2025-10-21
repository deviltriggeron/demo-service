package ports

import (
	"context"

	e "demo-service/internal/entity"
)

type Broker interface {
	CreateTopic(topic string, numPartitions, replicationFactor int)
	ConsumeOrder(ctx context.Context, topic string, handler func(order e.Order, method string) error) error
	ProduceOrder(topic string, method string, order e.Order) error
}
