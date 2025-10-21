package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"demo-service/internal/entity"
	"demo-service/internal/infrastructure/mocks"
	"demo-service/internal/service"
)

func TestService_Storage_Get(t *testing.T) {
	cacheMock := &mocks.OrderCache{}
	storageMock := &mocks.OrderRepository{}

	order := entity.Order{OrderUID: "123"}

	cacheMock.On("Get", "123").Return(entity.Order{}, false)
	storageMock.On("MethodSelect", "123").Return(order, nil)
	cacheMock.On("Set", order).Return()

	svc := service.NewOrderService(cacheMock, storageMock)

	got, err := svc.Get("123")
	assert.NoError(t, err)
	assert.Equal(t, order, got)

	cacheMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
}

func TestService_Storage_GetAll(t *testing.T) {
	cacheMock := &mocks.OrderCache{}
	storageMock := &mocks.OrderRepository{}

	var orders []entity.Order

	storageMock.On("MethodSelectAll").Return(orders, nil)
	cacheMock.On("SetAll", orders).Return()

	svc := service.NewOrderService(cacheMock, storageMock)

	got, err := svc.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, orders, got)

	storageMock.AssertExpectations(t)
}

func TestService_Storage_Set(t *testing.T) {
	cacheMock := &mocks.OrderCache{}
	storageMock := &mocks.OrderRepository{}

	order := returnOrder()

	storageMock.On("MethodInsert", order).Return(nil)
	cacheMock.On("Set", order).Return()

	svc := service.NewOrderService(cacheMock, storageMock)

	err := svc.InsertOrder(order)
	assert.NoError(t, err)

	storageMock.AssertExpectations(t)
}

func TestService_Storage_SetAll(t *testing.T) {
	cacheMock := &mocks.OrderCache{}
	storageMock := &mocks.OrderRepository{}

	var orders []entity.Order
	orders = append(orders, returnOrder())

	for i := range orders {
		storageMock.On("MethodInsert", orders[i]).Return(nil)
	}
	cacheMock.On("SetAll", orders).Return()

	svc := service.NewOrderService(cacheMock, storageMock)
	err := svc.InsertOrderAll(orders)
	assert.NoError(t, err)

	storageMock.AssertExpectations(t)
}

func TestService_Storage_Update(t *testing.T) {
	cacheMock := &mocks.OrderCache{}
	storageMock := &mocks.OrderRepository{}

	order := returnOrder()

	cacheMock.On("Delete", "a12345b6c7890test").Return(false)
	cacheMock.On("Set", order).Return(false)
	storageMock.On("MethodUpdate", order).Return(nil)

	svc := service.NewOrderService(cacheMock, storageMock)

	err := svc.UpdateOrder(order)
	assert.NoError(t, err)

	storageMock.AssertExpectations(t)
}

func TestService_Storage_Delete(t *testing.T) {
	cacheMock := &mocks.OrderCache{}
	storageMock := &mocks.OrderRepository{}

	order := returnOrder()

	cacheMock.On("Delete", order.OrderUID).Return(true)
	storageMock.On("MethodDelete", order.OrderUID).Return(nil)

	svc := service.NewOrderService(cacheMock, storageMock)

	err := svc.DeleteOrder(order.OrderUID)
	assert.NoError(t, err)

	storageMock.AssertExpectations(t)
}
