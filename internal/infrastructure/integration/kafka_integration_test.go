package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"demo-service/internal/entity"
	"demo-service/internal/infrastructure/mocks"
)

func TestService_CreateTopic(t *testing.T) {
	brokerMock := &mocks.Broker{}

	brokerMock.On("CreateTopic", "orders", 1, 1).Return()

	brokerMock.CreateTopic("orders", 1, 1)

	brokerMock.AssertExpectations(t)
}

func TestService_SendOrderToBroker(t *testing.T) {
	brokerMock := &mocks.Broker{}
	order := entity.Order{OrderUID: "123"}
	topic := "orders"
	method := "create"

	brokerMock.On("ProduceOrder", topic, method, order).Return(nil)

	err := brokerMock.ProduceOrder(topic, method, order)

	assert.NoError(t, err)

	brokerMock.AssertExpectations(t)
}

func TestService_HandleKafkaMessage(t *testing.T) {
	brokerMock := &mocks.Broker{}
	order := entity.Order{OrderUID: "123"}
	topic := "orders"

	brokerMock.On("ConsumeOrder", mock.Anything, topic, mock.AnythingOfType("func(entity.Order, string) error")).
		Return(func(ctx context.Context, t string, handler func(entity.Order, string) error) error {
			return handler(order, "create")
		})

	handler := func(o entity.Order, method string) error {
		assert.Equal(t, order, o)
		assert.Equal(t, "create", method)
		return nil
	}

	err := brokerMock.ConsumeOrder(context.Background(), topic, handler)
	assert.NoError(t, err)

	brokerMock.AssertExpectations(t)
}
