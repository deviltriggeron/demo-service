package service

import (
	"context"
	"errors"
	"fmt"

	e "demo-service/internal/entity"
	ports "demo-service/internal/interface"
)

type orderService struct {
	cache ports.OrderCache
	db    ports.OrderRepository
}

func NewOrderService(c ports.OrderCache, storage ports.OrderRepository) OrderService {
	return &orderService{
		cache: c,
		db:    storage,
	}
}

func (s *orderService) Get(id string) (e.Order, error) {
	if order, ok := s.cache.Get(id); ok {
		return order, nil
	}

	order, err := s.db.MethodSelect(id)
	if err != nil {
		return e.Order{}, err
	}

	s.cache.Set(order)
	return order, nil
}

func (s *orderService) GetAll() ([]e.Order, error) {
	orders, err := s.db.MethodSelectAll()
	if err != nil {
		return []e.Order{}, err
	}

	s.cache.SetAll(orders)
	return orders, nil
}

func (s *orderService) InsertOrder(order e.Order) error {
	if err := validateOrder(order); err != nil {
		return err
	}
	if err := s.db.MethodInsert(order); err != nil {
		return err
	}
	s.cache.Set(order)
	return nil
}

func (s *orderService) InsertOrderAll(orders []e.Order) error {
	for i := range orders {
		if err := validateOrder(orders[i]); err != nil {
			return err
		}
		if err := s.db.MethodInsert(orders[i]); err != nil {
			return err
		}
	}
	s.cache.SetAll(orders)
	return nil
}

func (s *orderService) DeleteOrder(id string) error {
	if err := s.db.MethodDelete(id); err != nil {
		return err
	}

	if !s.cache.Delete(id) {
		return fmt.Errorf("error delete in cache")
	}

	return nil
}

func (s *orderService) UpdateOrder(order e.Order) error {
	if err := s.db.MethodUpdate(order); err != nil {
		return err
	}
	s.cache.Delete(order.OrderUID)
	s.cache.Set(order)
	return nil
}

func (s *orderService) Shutdown(ctx context.Context) error {
	return s.db.CloseDB()
}

func (s *orderService) HandleKafkaMessage(order e.Order, method string) error {
	switch method {
	case "INSERT":
		return s.InsertOrder(order)
	case "UPDATE":
		return s.UpdateOrder(order)
	case "DELETE":
		return s.DeleteOrder(order.OrderUID)
	}
	return fmt.Errorf("unknown method: %s", method)
}

func validateOrder(o e.Order) error {
	if o.OrderUID == "" {
		return fmt.Errorf("order_uid is required")
	}
	if o.TrackNumber == "" {
		return errors.New("track_number is required")
	}
	if o.DeliveryService == "" {
		return fmt.Errorf("delivery_service is required")
	}
	if o.Delivery.Address == "" {
		return fmt.Errorf("delivery.address is required")
	}
	for _, item := range o.Items {
		if item.Price <= 0 {
			return fmt.Errorf("item price must be positive")
		}
	}
	return nil
}
