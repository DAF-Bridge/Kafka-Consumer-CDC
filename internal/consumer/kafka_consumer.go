package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/DAF-Bridge/cdc-service/internal/models"
	"github.com/DAF-Bridge/cdc-service/internal/repository"
	"github.com/DAF-Bridge/cdc-service/internal/service"
	"github.com/DAF-Bridge/cdc-service/pkg/logs"
	"github.com/segmentio/kafka-go"
)

type OpenSearchCDC struct {
	repo    repository.OpenSearchRepository
	service service.OpenSearchService
}

func NewOpenSearchCDC(repo repository.OpenSearchRepository, service service.OpenSearchService) *OpenSearchCDC {
	return &OpenSearchCDC{
		repo:    repo,
		service: service,
	}
}

func (opn *OpenSearchCDC) ConsumeMessage(message kafka.Message) error {
	var event models.CDCEvent
	if event.Payload.Op == "r" {
		logs.Info(fmt.Sprintf("Ignoring read operation (r) for event: %v", event))
		return nil
	}

	if event.Payload.Op == "d" {
		err := opn.service.ProcessEvent(event)
		if err != nil {
			logs.Error(fmt.Sprintf("Error processing event: %v", err))
			return fmt.Errorf("error processing event: %v", err)
		}
	}

	if err := json.Unmarshal(message.Value, &event); err != nil {
		return fmt.Errorf("error unmarshalling event: %v", err)
	}
	// Process only Create, Update, Delete events
	switch event.Payload.Op {
	case "c", "u", "d":
		// Delegate processing to the service layer
		err := opn.service.ProcessEvent(event)
		if err != nil {
			logs.Error(fmt.Sprintf("Error processing event: %v", err))
			return fmt.Errorf("error processing event: %v", err)
		}
	default:
		logs.Info(fmt.Sprintf("Ignoring read operation (r) for event: %v", event))
		// Ignore "r" (read) operations
		return nil
	}

	return nil
}

func StartKafka(cdcConsumer *OpenSearchCDC) {
	broker := os.Getenv("KAFKA_BROKER")
	topic := os.Getenv("KAFKA_TOPIC")
	groupID := string(os.Getenv("KAFKA_GROUP_ID"))

	config := kafka.ReaderConfig{
		Brokers:  []string{broker},
		GroupID:  groupID,
		Topic:    topic,
		MaxBytes: 10e6, // 10MB
	}

	reader := kafka.NewReader(config)

	for {
		msg, err := reader.ReadMessage(context.Background())

		if err != nil {
			log.Fatal(err)
		}

		logs.Info(fmt.Sprintf("Message received: %v", string(msg.Value)))

		if err := cdcConsumer.ConsumeMessage(msg); err != nil {
			logs.Error(fmt.Sprintf("Error consuming message: %v", err))
		}
	}
}
