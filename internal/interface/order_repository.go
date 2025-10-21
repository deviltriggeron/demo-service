package ports

import e "demo-service/internal/entity"

type OrderRepository interface {
	CloseDB() error
	MethodSelect(id string) (e.Order, error)
	MethodSelectAll() ([]e.Order, error)
	MethodDelete(orderUID string) error
	MethodInsert(orders e.Order) error
	MethodUpdate(order e.Order) error
}
