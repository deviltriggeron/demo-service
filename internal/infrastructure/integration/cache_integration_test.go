package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"demo-service/internal/entity"
	"demo-service/internal/infrastructure/mocks"
	"demo-service/internal/service"
)

func TestService_Cache_Get(t *testing.T) {
	cacheMock := &mocks.OrderCache{}
	storageMock := &mocks.OrderRepository{}

	order := entity.Order{OrderUID: "123"}

	cacheMock.On("Get", "123").Return(order, true)

	svc := service.NewOrderService(cacheMock, storageMock)

	got, err := svc.Get("123")
	assert.NoError(t, err)
	assert.Equal(t, order, got)

	cacheMock.AssertExpectations(t)
}

func TestService_Cache_Set(t *testing.T) {
	cacheMock := &mocks.OrderCache{}
	storageMock := &mocks.OrderRepository{}

	order := returnOrder()

	storageMock.On("MethodInsert", order).Return(nil)
	cacheMock.On("Set", order).Return()

	svc := service.NewOrderService(cacheMock, storageMock)

	err := svc.InsertOrder(order)
	assert.NoError(t, err)

	cacheMock.AssertExpectations(t)
}

func TestService_Cache_SetAll(t *testing.T) {
	cacheMock := &mocks.OrderCache{}
	storageMock := &mocks.OrderRepository{}

	var orders []entity.Order
	orders = append(orders, returnOrder())

	cacheMock.On("SetAll", orders).Return()
	for i := range orders {
		storageMock.On("MethodInsert", orders[i]).Return(nil)
	}

	svc := service.NewOrderService(cacheMock, storageMock)

	err := svc.InsertOrderAll(orders)
	assert.NoError(t, err)

	cacheMock.AssertExpectations(t)
}

func TestService_Cache_Delete(t *testing.T) {
	cacheMock := &mocks.OrderCache{}
	storageMock := &mocks.OrderRepository{}

	cacheMock.On("Delete", "a12345b6c7890test").Return(true)

	svc := service.NewOrderService(cacheMock, storageMock)

	err := svc.DeleteOrder("a12345b6c7890test")
	assert.NoError(t, err)

	cacheMock.AssertExpectations(t)
}

func TestService_Cache_Update(t *testing.T) {
	cacheMock := &mocks.OrderCache{}
	storageMock := &mocks.OrderRepository{}

	order := returnOrder()

	cacheMock.On("Delete", "a12345b6c7890test").Return(true)
	cacheMock.On("Set", order).Return(true)

	svc := service.NewOrderService(cacheMock, storageMock)

	err := svc.UpdateOrder(order)
	assert.NoError(t, err)

	cacheMock.AssertExpectations(t)
}

func returnOrder() entity.Order {
	return entity.Order{
		OrderUID:        "a12345b6c7890test",
		TrackNumber:     "WBNEWTRACK2025",
		Entry:           "WBX",
		Locale:          "ru",
		InternalSig:     "sig123",
		CustomerID:      "customer42",
		DeliveryService: "dhl",
		Shardkey:        "5",
		SmID:            77,
		DateCreated:     "2025-08-27T12:45:00Z",
		OofShard:        "3",

		Delivery: entity.Delivery{
			Name:    "Ivan Ivanov",
			Phone:   "+79001234567",
			Zip:     "101000",
			City:    "Moscow",
			Address: "Tverskaya 10",
			Region:  "Moscow Region",
			Email:   "ivanov@example.com",
		},

		Payment: entity.Payment{
			Transaction:  "a12345b6c7890test",
			RequestID:    "req-555",
			Currency:     "EUR",
			Provider:     "paypal",
			Amount:       2599,
			PaymentDT:    1693123456,
			Bank:         "sberbank",
			DeliveryCost: 499,
			GoodsTotal:   2100,
			CustomFee:    0,
		},

		Items: []entity.Item{
			{
				ChrtID:      11223344,
				TrackNumber: "WBNEWTRACK2025",
				Price:       700,
				Rid:         "rid-xyz-001",
				Name:        "Wireless Headphones",
				Sale:        15,
				Size:        "M",
				TotalPrice:  595,
				NmID:        445566,
				Brand:       "Sony",
				Status:      201,
			},
		},
	}
}
