package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"

	ports "demo-service/internal/interface"
)

type kafkaBroker struct {
	brokers []string
	groupID string
}

func NewKafkaBroker(brokers []string, groupID string) ports.Broker {
	return &kafkaBroker{
		brokers: brokers,
		groupID: groupID,
	}
}

func (k *kafkaBroker) CreateTopic(topic string, numPartitions, replicationFactor int) {
	if err := createKafkaTopic(topic, numPartitions, replicationFactor); err != nil {
		log.Fatal(err)
	}
	if err := createKafkaTopic(topic+".DLQ", numPartitions, replicationFactor); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Topic %q and DLQ created\n", topic)
}

func createKafkaTopic(topic string, numPartitions, replicationFactor int) error {
	conn, err := kafka.Dial("tcp", "kafka:9092")
	if err != nil {
		return fmt.Errorf("failed to dial kafka: %w", err)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("error close connection kafka: %v", err)
		}
	}()

	err = conn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     numPartitions,
		ReplicationFactor: replicationFactor,
	})
	if err != nil {
		return fmt.Errorf("failed to create topic %s: %w", topic, err)
	}
	return nil
}

func (k *kafkaBroker) ClearTopic(topic string, partition int) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", "kafka:9092", topic, partition)
	if err != nil {
		log.Fatalf("failed to dial leader: %v", err)
	}

	err = conn.DeleteTopics(topic)
	if err != nil {
		log.Fatalf("failed delete topics: %v", err)
	}
}
