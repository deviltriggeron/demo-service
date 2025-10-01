package cache

import (
	model "demo-service/internal/model"
	"sync"
)

type OrderCache struct {
	mu     sync.Mutex
	orders map[string]model.Order
}

func NewCache() *OrderCache {
	return &OrderCache{
		orders: make(map[string]model.Order),
	}
}

func (c *OrderCache) Get(id string) (model.Order, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	o, ok := c.orders[id]
	return o, ok
}

func (c *OrderCache) Set(order model.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.orders[order.OrderUID] = order
}

func (c *OrderCache) SetAll(orders []model.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i := range orders {
		c.orders[orders[i].OrderUID] = orders[i]
	}
}

func (c *OrderCache) Delete(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.orders, id)
}

func (c *OrderCache) DeleteAll() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k := range c.orders {
		delete(c.orders, k)
	}
}

func (c *OrderCache) GetAll() []model.Order {
	c.mu.Lock()
	defer c.mu.Unlock()

	orders := make([]model.Order, 0, len(c.orders))

	for _, o := range c.orders {
		orders = append(orders, o)
	}
	return orders
}
