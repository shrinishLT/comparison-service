package config

import (
	"os"
	"smartui-comparison-service/constants"
	"strconv"
	"strings"
)

type Config struct {
	KafkaBrokers    []string
	ComparisonTopic string
	ResultTopic     string
	MaxWorkers      int
}

func LoadConfig() Config {
	maxWorkers, err := strconv.Atoi(os.Getenv("MAX_WORKERS"))
	if err != nil {
		maxWorkers = 10
	}

	return Config{
		KafkaBrokers:    strings.Split(constants.KAFKA_BROKERS, ","),
		ComparisonTopic: os.Getenv("COMPARISON_TOPIC"),
		ResultTopic:     os.Getenv("RESULT_TOPIC"),
		MaxWorkers:      maxWorkers,
	}
}
