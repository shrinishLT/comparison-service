package constants

import "time"

const (
	KAFKA_BROKERS = "localhost:9092"
	MAX_WORKERS   = 10
)

const (
	WAIT_GROUP_TIMEOUT = 300 * time.Second
)
