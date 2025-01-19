package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"smartui-comparison-service/config"
	"smartui-comparison-service/constants"
	"smartui-comparison-service/kafka"
	"smartui-comparison-service/models"
	"smartui-comparison-service/waitGroup"
	"smartui-comparison-service/workers"
	"syscall"
)

func main() {

	cfg := config.LoadConfig()

	producer := kafka.NewProducer(cfg.KafkaBrokers, cfg.ResultTopic)
	defer func() {
		if err := producer.Close(); err != nil {
			log.Printf("Error closing Kafka producer: %v", err)
		} else {
			log.Println("Kafka producer shut down gracefully.")
		}
	}()

	resultHandler := func(ctx context.Context, result models.ComparisonResult) error {
		return producer.PublishResult(ctx, result)
	}

	workerPool := workers.NewWorkerPool(cfg.MaxWorkers, resultHandler)

	consumer := kafka.NewConsumer(cfg.KafkaBrokers, cfg.ComparisonTopic, "comparison-consumer-group")
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Printf("Error closing Kafka consumer: %v", err)
		} else {
			log.Println("Kafka consumer shut down gracefully.")
		}
	}()

	// Create a context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())

	// Signal handling for shutdown
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		log.Println("Received shutdown signal...")
		cancel()
	}()

	// Start Kafka consumer in a goroutine
	log.Println("Starting Kafka consumer...")
	go func() {
		consumer.ConsumeMessages(ctx, workerPool)
	}()

	// Block until context is canceled
	<-ctx.Done()
	log.Println("Shutting down services...")

	// Wait for worker pool tasks to complete
	if waitGroup.WaitWithTimeout(waitGroup.GetGlobalWaitGroup(), constants.WAIT_GROUP_TIMEOUT) {
		log.Println("all goroutines are finished")
	} else {
		log.Println("all goroutines weren't able to complete on time. Initiating cleanup of resources")
	}

	log.Println("All worker tasks completed.")

	if err := kafka.Close(); err != nil {
		logger.Errorf("Kafka cleanup error %s", err.Error())
	}
}
