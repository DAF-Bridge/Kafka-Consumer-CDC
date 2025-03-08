package initializer

import (
	"fmt"
	"log"
	"os"
	"strings"
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
	// consumer, err := kafka.NewConsumer(config)
	// if err != nil {
	// 	log.Fatalf("Failed to create Kafka consumer: %v", err)
	// }
	// defer consumer.Close()

	// // Check metadata to verify broker connection
	// _, err = consumer.GetMetadata(nil, true, int(5*time.Second.Milliseconds()))
	// if err != nil {
	// 	log.Fatalf("Failed to get Kafka metadata: %v", err)
	// }
	topicsStr := os.Getenv("KAFKA_TOPIC")
	topics := strings.Split(topicsStr, ",")
	// topics := os.Getenv("KAFKA_TOPIC")

	// // Check Kafka metadata for the specified topic
	// meta, err := consumer.GetMetadata(&topic, false, int(5*time.Second.Milliseconds()))
	// if err != nil {
	// 	log.Fatalf("Failed to get Kafka metadata for topic %s: %v", topic, err)
	// }
	// topicMeta, exists := meta.Topics[topic]
	// if !exists || len(topicMeta.Partitions) == 0 {
	// 	log.Fatalf("Topic %s exists but has no partitions", topic)
	// }

	// logs.Info(fmt.Sprintf("Successfully connected to Kafka topic: %s", topic))

	// Create a Kafka consumer
	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer consumer.Close()

	// Check Kafka metadata to verify broker connection
	meta, err := consumer.GetMetadata(nil, true, int(5*time.Second.Milliseconds()))
	if err != nil {
		log.Fatalf("Failed to get Kafka metadata: %v", err)
	}

	// Check if all required topics exist
	for _, topic := range topics {
		topicMeta, exists := meta.Topics[topic]
		if !exists {
			log.Fatalf("Topic %s does not exist in Kafka", topic)
		}
		if len(topicMeta.Partitions) == 0 {
			log.Fatalf("Topic %s exists but has no partitions", topic)
		}

		logs.Info(fmt.Sprintf("Successfully connected to Kafka topic: %s", topic))
	}

	logs.Info("Kafka connection check completed successfully.")
}
