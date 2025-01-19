package workers

import (
	"context"
	"fmt"
	"smartui-comparison-service/comparison"
	"smartui-comparison-service/models"
)

type WorkerPool struct {
	sem           chan struct{}
	resultHandler func(ctx context.Context, result models.ComparisonResult) error
}

func NewWorkerPool(maxWorkers int, resultHandler func(ctx context.Context, result models.ComparisonResult) error) *WorkerPool {
	return &WorkerPool{
		sem:           make(chan struct{}, maxWorkers),
		resultHandler: resultHandler,
	}
}

func (p *WorkerPool) Enqueue(ctx context.Context, task models.ComparisonTask) {
	p.sem <- struct{}{} // acquire semaphore slot
	go func() {
		defer func() { <-p.sem }() // release semaphore slot

		// process the task
		result := comparison.ProcessTask(task)

		// handle the result ( publish it to Kafka)
		if err := p.resultHandler(ctx, result); err != nil {
			fmt.Printf("Failed to handle result for %s: %v\n", task.RequestID, err)
		}
	}()
}
