package ports

import e "demo-service/internal/entity"

type OrderCache interface {
	Get(id string) (e.Order, bool)
	Set(order e.Order)
	SetAll(orders []e.Order)
	Delete(id string) bool
	DeleteAll()
	GetAll() []e.Order
}
