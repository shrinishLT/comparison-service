package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"smartui-comparison-service/models"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
	return &Producer{writer: writer}
}

func (p *Producer) PublishResult(ctx context.Context, result models.ComparisonResult) error {
	messageBytes, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to serialize result: %w", err)
	}

	if err := p.writer.WriteMessages(ctx, kafka.Message{
		Value: messageBytes,
	}); err != nil {
		log.Printf("Failed to publish message to Kafka: %v", err)
		return err
	}

	log.Printf("Published result to Kafka: %s", string(messageBytes))
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
