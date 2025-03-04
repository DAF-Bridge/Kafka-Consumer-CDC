package consumer

import (
	"encoding/json"
	"fmt"

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
