package kafka

import (
	"context"
	"encoding/json"
	"log"
	"smartui-comparison-service/models"
	"smartui-comparison-service/workers"
	"time"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, topic, groupID string) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        groupID,
		MaxBytes:       10e6, // 10MB
		CommitInterval: time.Second,
	})
	return &Consumer{reader: reader}
}

func (c *Consumer) ConsumeMessages(ctx context.Context, workerPool *workers.WorkerPool) {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}

		var task models.ComparisonTask
		if err := json.Unmarshal(msg.Value, &task); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		log.Printf("Received task: %v", task.RequestID)
		workerPool.Enqueue(ctx, task)
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
