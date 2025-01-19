package kafka_test

import (
	"context"
	"smartui-comparison-service/kafka"
	"smartui-comparison-service/models"
	"testing"
	"time"

	"github.com/segmentio/kafka-go/kafkatest"
)

func TestKafkaProducer(t *testing.T) {
	// Create a mock Kafka writer
	mockWriter := &kafkatest.Writer{}

	// Initialize the producer
	producer := &kafka.Producer{
		Writer: mockWriter,
	}

	// Create a test message
	result := models.ComparisonResult{
		RequestID: "test-request",
		Success:   true,
		Message:   "Task processed successfully",
	}

	// Publish the result
	err := producer.PublishResult(context.Background(), result)
	if err != nil {
		t.Fatalf("failed to publish message: %v", err)
	}

	// Verify the message was written
	if len(mockWriter.Messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(mockWriter.Messages))
	}

	// Assert the message content
	expectedMessage := `{"RequestID":"test-request","Success":true,"Message":"Task processed successfully"}`
	actualMessage := string(mockWriter.Messages[0].Value)
	if actualMessage != expectedMessage {
		t.Errorf("unexpected message content: got %v, want %v", actualMessage, expectedMessage)
	}
}

func TestKafkaConsumer(t *testing.T) {
	// Create a mock Kafka reader
	mockReader := &kafkatest.Reader{
		Messages: []kafka.Message{
			{Value: []byte(`{"RequestID":"test-task1"}`)},
			{Value: []byte(`{"RequestID":"test-task2"}`)},
		},
	}

	// Initialize the consumer
	consumer := &kafka.Consumer{
		Reader: mockReader,
	}

	// Create a mock worker pool
	var tasks []models.ComparisonTask
	workerPool := &workers.MockWorkerPool{
		EnqueueFn: func(ctx context.Context, task interface{}) {
			tasks = append(tasks, task.(models.ComparisonTask))
		},
	}

	// Consume messages
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go consumer.ConsumeMessages(ctx, workerPool)

	// Wait for the consumer to process the messages
	time.Sleep(1 * time.Second)

	// Assert that all tasks were enqueued
	if len(tasks) != len(mockReader.Messages) {
		t.Fatalf("expected %d tasks, got %d", len(mockReader.Messages), len(tasks))
	}

	// Validate the task content
	expectedTasks := []models.ComparisonTask{
		{RequestID: "test-task1"},
		{RequestID: "test-task2"},
	}
	for i, task := range tasks {
		if task.RequestID != expectedTasks[i].RequestID {
			t.Errorf("unexpected task content: got %v, want %v", task.RequestID, expectedTasks[i].RequestID)
		}
	}
}
