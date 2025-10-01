package kafka

import (
	"context"
	"encoding/json"
	"log"

	model "demo-service/internal/model"

	"github.com/segmentio/kafka-go"
)

func ProduceOrder(topic string, method string, order model.Order) error {
	w := &kafka.Writer{
		Addr:     kafka.TCP("kafka:9092"),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	defer func() {
		if err := w.Close(); err != nil {
			log.Fatal("failed to close writer:", err)
		}
	}()

	msg := model.KafkaMsg{
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
		log.Fatal("failed to write message:", err)
	}

	return nil
}
