package service

import (
	"context"

	e "demo-service/internal/entity"
)

type OrderService interface {
	Get(id string) (e.Order, error)
	GetAll() ([]e.Order, error)
	InsertOrder(order e.Order) error
	InsertOrderAll(orders []e.Order) error
	DeleteOrder(id string) error
	UpdateOrder(order e.Order) error
	Shutdown(ctx context.Context) error
	HandleKafkaMessage(order e.Order, method string) error
}
