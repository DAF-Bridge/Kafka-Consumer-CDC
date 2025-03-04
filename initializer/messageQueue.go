package initializer

import (
	"log"
	"os"
	"time"

	"github.com/DAF-Bridge/cdc-service/pkg/logs"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// ConnectToMessageQueue checks the health of the Kafka connection for a consumer
func ConnectToMessageQueue() {
	broker := os.Getenv("KAFKA_BROKER") // Example: "kafka:9092"
	if broker == "" {
		broker = "localhost:9092" // Default fallback
	}

	config := &kafka.ConfigMap{
		"bootstrap.servers":  broker,
		"group.id":           "health-check-consumer", // Unique group for health checks
		"auto.offset.reset":  "earliest",              // Start from the beginning if no offsets exist
		"enable.auto.commit": false,                   // Disable auto-commit for health checks
	}

	// Create a Kafka consumer
	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer consumer.Close()

	// Check metadata to verify broker connection
	_, err = consumer.GetMetadata(nil, true, int(5*time.Second.Milliseconds()))
	if err != nil {
		log.Fatalf("Failed to get Kafka metadata: %v", err)
	}

	logs.Info("Successfully connected to Kafka as a consumer!")
}
