package cache

import (
	"time"

	"github.com/patrickmn/go-cache"

	e "demo-service/internal/entity"
	p "demo-service/internal/interface"
)

type orderCache struct {
	cache *cache.Cache
}

func NewCache(defaultExpiration, cleanupInterval time.Duration) p.OrderCache {
	c := cache.New(defaultExpiration, cleanupInterval)
	return &orderCache{cache: c}
}

func (c *orderCache) Get(id string) (e.Order, bool) {
	v, found := c.cache.Get(id)
	if !found {
		return e.Order{}, false
	}
	return v.(e.Order), true
}

func (c *orderCache) Set(order e.Order) {
	c.cache.Set(order.OrderUID, order, cache.DefaultExpiration)
}

func (c *orderCache) SetAll(orders []e.Order) {
	for _, order := range orders {
		c.cache.Set(order.OrderUID, order, cache.DefaultExpiration)
	}
}

func (c *orderCache) Delete(id string) bool {
	c.cache.Delete(id)
	return true
}

func (c *orderCache) DeleteAll() {
	for k := range c.cache.Items() {
		c.cache.Delete(k)
	}
}

func (c *orderCache) GetAll() []e.Order {
	items := c.cache.Items()
	orders := make([]e.Order, 0, len(items))
	for _, item := range items {
		if order, ok := item.Object.(e.Order); ok {
			orders = append(orders, order)
		}
	}
	return orders
}
