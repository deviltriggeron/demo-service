package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func CreateTopic(topic string, numPartitions int, replicationFactor int) {
	conn, err := kafka.Dial("tcp", "localhost:9092")
	if err != nil {
		log.Fatal("failed to dial kafka:", err)
	}
	defer conn.Close()

	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = conn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     numPartitions,
		ReplicationFactor: replicationFactor,
	})
	if err != nil {
		log.Fatalf("failed to create topic %s: %v", topic, err)
	}

	fmt.Printf("Topic %q created (partitions=%d, replication=%d)\n",
		topic, numPartitions, replicationFactor)
}

func ClearTopic(topic string, partition int) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	conn.DeleteTopics(topic)
}
