package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"

	e "demo-service/internal/entity"
)

func (k *kafkaBroker) ProduceOrder(topic string, method string, order e.Order) error {
	w := &kafka.Writer{
		Addr:     kafka.TCP("kafka:9092"),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	defer func() {
		if err := w.Close(); err != nil {
			log.Printf("failed to close writer: %v", err)
		}
	}()

	msg := e.KafkaMsg{
		Order:  order,
		Method: method,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = w.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte("order"),
			Value: data,
		},
	)
	if err != nil {
		log.Printf("failed to write message: %v", err)
	}

	return nil
}
