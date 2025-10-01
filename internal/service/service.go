package service

import (
	"context"
	cache "demo-service/internal/cache"
	model "demo-service/internal/model"
	postgres "demo-service/internal/postgres"
)

type OrderService struct {
	cache *cache.OrderCache
}

func NewOrderService(c *cache.OrderCache) *OrderService {
	return &OrderService{cache: c}
}

func (s OrderService) Get(id string) (model.Order, error) {
	if order, ok := s.cache.Get(id); ok {
		return order, nil
	}

	order, err := postgres.MethodSelect(id)
	if err != nil {
		return model.Order{}, err
	}

	return order, nil
}

func (s OrderService) GetAll() ([]model.Order, error) {
	orders := s.cache.GetAll()
	if len(orders) > 0 {
		return orders, nil
	}
	orders, err := postgres.MethodSelectAll()
	if err != nil {
		return []model.Order{}, err
	}
	s.cache.SetAll(orders)
	return orders, nil
}

func (s *OrderService) InsertOrder(order model.Order) error {
	if err := postgres.MethodInsert(order); err != nil {
		return err
	}
	s.cache.Set(order)
	return nil
}

func (s *OrderService) InsertOrderAll(orders []model.Order) error {
	for i := range orders {
		if err := postgres.MethodInsert(orders[i]); err != nil {
			return err
		}
	}
	s.cache.SetAll(orders)
	return nil
}

func (s *OrderService) DeleteOrder(id string) error {
	if err := postgres.MethodDelete(id); err != nil {
		return err
	}
	s.cache.Delete(id)
	return nil
}

func (s *OrderService) UpdateOrder(order model.Order) error {
	if err := postgres.MethodUpdate(order); err != nil {
		return err
	}
	s.cache.Delete(order.OrderUID)
	s.cache.Set(order)
	return nil
}

func (s OrderService) Shutdown(ctx context.Context) error {
	return postgres.CloseDB()
}
