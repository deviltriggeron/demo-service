package kafka

import (
	"context"
	model "demo-service/internal/model"
	service "demo-service/internal/service"
	"encoding/json"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

func ConsumeOrder(ctx context.Context, topic string, svc *service.OrderService) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{"kafka:9092"},
		Topic:       topic,
		Partition:   0,
		MinBytes:    10e3,
		MaxBytes:    10e6,
		StartOffset: kafka.FirstOffset,
	})

	defer func() {
		if err := r.Close(); err != nil {
			log.Fatal("failed to close reader:", err)
		}
	}()

	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			log.Printf("error while reading message: %v", err)
			break
		}

		var msg model.KafkaMsg
		if err := json.Unmarshal(m.Value, &msg); err != nil {
			log.Println("Error parsing JSON:", err)
			continue
		}

		switch msg.Method {
		case "UPDATE":
			svc.UpdateOrder(msg.Order)
		case "DELETE":
			svc.DeleteOrder(msg.Order.OrderUID)
		case "INSERT":
			svc.InsertOrder(msg.Order)
		default:
			fmt.Println("Unknown action: ", msg.Method)
		}
	}
}
