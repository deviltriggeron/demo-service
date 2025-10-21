package cache

import (
	lru "github.com/hashicorp/golang-lru/v2"

	e "demo-service/internal/entity"
	p "demo-service/internal/interface"
)

type orderCache struct {
	cache *lru.Cache[string, e.Order]
}

func NewCache(size int) (p.OrderCache, error) {
	c, err := lru.New[string, e.Order](size)
	if err != nil {
		return nil, err
	}
	return &orderCache{cache: c}, nil
}

func (c *orderCache) Get(id string) (e.Order, bool) {
	return c.cache.Get(id)
}

func (c *orderCache) GetAll() []e.Order {
	orders := make([]e.Order, 0, len(c.cache.Values()))
	orders = append(orders, c.cache.Values()...)

	return orders
}

func (c *orderCache) Set(order e.Order) {
	c.cache.Add(order.OrderUID, order)
}

func (c *orderCache) SetAll(orders []e.Order) {
	for i := range orders {
		c.cache.Add(orders[i].OrderUID, orders[i])
	}
}

func (c *orderCache) Delete(id string) bool {
	return c.cache.Remove(id)
}

func (c *orderCache) DeleteAll() {
	for _, key := range c.cache.Keys() {
		c.cache.Remove(key)
	}
}
