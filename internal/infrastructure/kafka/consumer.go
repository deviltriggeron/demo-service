package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/segmentio/kafka-go"

	e "demo-service/internal/entity"
)

func (k *kafkaBroker) ConsumeOrder(ctx context.Context, topic string, handler func(order e.Order, method string) error) error {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: k.brokers,
		Topic:   topic,
		GroupID: "order-consumer",
	})
	defer func() {
		if err := r.Close(); err != nil {
			log.Printf("failed to close Kafka reader: %v", err)
		}
	}()

	const maxRetries = 3

	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			return err
		}

		var msg e.KafkaMsg
		if err := json.Unmarshal(m.Value, &msg); err != nil {
			log.Printf("failed to unmarshal message: %v", err)
			continue
		}

		var retryErr error
		for i := 0; i < maxRetries; i++ {
			retryErr = handler(msg.Order, msg.Method)
			if errors.Is(retryErr, nil) {
				break
			}
			log.Printf("retry %d/%d failed for order %s: %v", i+1, maxRetries, msg.Order.OrderUID, retryErr)
			sleep := time.Second * time.Duration(1<<i)
			time.Sleep(sleep)
		}

		if retryErr != nil {
			if err := k.ProduceOrder(topic+".DLQ", msg.Method, msg.Order); err != nil {
				log.Printf("failed to send to DLQ: %v", err)
			}
		}
	}
}
