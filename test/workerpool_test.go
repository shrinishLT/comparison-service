package workers_test

import (
	"context"
	"smartui-comparison-service/workers"
	"testing"
	"time"
)

func TestWorkerPool(t *testing.T) {

	var results []string
	resultHandler := func(ctx context.Context, result interface{}) error {
		results = append(results, result.(string))
		return nil
	}

	pool := workers.NewWorkerPool(2, resultHandler)

	// Enqueue tasks
	tasks := []string{"task1", "task2", "task3"}
	for _, task := range tasks {
		pool.Enqueue(context.Background(), task)
	}

	// Wait for all tasks to complete
	time.Sleep(10 * time.Second)

	// Assert results
	if len(results) != len(tasks) {
		t.Fatalf("expected %d results, got %d", len(tasks), len(results))
	}
	for i, task := range tasks {
		if results[i] != "Processed: "+task {
			t.Errorf("unexpected result: got %v, want %v", results[i], "Processed: "+task)
		}
	}
}
