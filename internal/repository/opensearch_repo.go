package repository

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/DAF-Bridge/cdc-service/errs"
	"github.com/DAF-Bridge/cdc-service/internal/models"
	"github.com/DAF-Bridge/cdc-service/pkg/logs"
	"github.com/opensearch-project/opensearch-go"
)

type openSearchRepository struct {
	es *opensearch.Client
}

func NewOpenSearchRepository(es *opensearch.Client) OpenSearchRepository {
	return &openSearchRepository{
		es: es,
	}
}

func (os *openSearchRepository) CreateOrUpdateEvent(event *models.EventDocument) error {
	// Convert event data to JSON
	data, err := json.Marshal(event)
	if err != nil {
		return errs.NewCannotBeProcessedError("error marshalling event data")
	}

	// Index the document
	res, err := os.es.Index("events", bytes.NewReader(data), os.es.Index.WithDocumentID(fmt.Sprintf("%d", event.ID)))
	if err != nil {
		return errs.NewCannotBeProcessedError("error indexing document")
	}
	defer res.Body.Close()

	logs.Info(fmt.Sprintf("Event document indexed: %v", event))

	return nil
}

func (os *openSearchRepository) DeleteEvent(event models.EventDocument) error {
	// Delete the document from OpenSearch
	res, err := os.es.Delete("events", fmt.Sprintf("%d", event.ID))
	if err != nil {
		return errs.NewCannotBeProcessedError(fmt.Sprintf("error deleting document: %v", err))
	}
	defer res.Body.Close()

	logs.Info(fmt.Sprintf("Event document deleted: %v", event))

	return nil
}

func (os *openSearchRepository) CreateOrUpdateJob(event *models.JobDocument) error {
	// Convert job data to JSON
	data, err := json.Marshal(event)
	if err != nil {
		return errs.NewCannotBeProcessedError("error marshalling job data")
	}

	// Index the document
	res, err := os.es.Index("jobs", bytes.NewReader(data), os.es.Index.WithDocumentID(fmt.Sprintf("%d", event.ID)))
	if err != nil {
		return errs.NewCannotBeProcessedError("error indexing document")
	}
	defer res.Body.Close()

	logs.Info(fmt.Sprintf("Job document indexed: %v", event))

	return nil
}

func (os *openSearchRepository) DeleteJob(event models.JobDocument) error {
	// Delete the document from OpenSearch
	res, err := os.es.Delete("jobs", fmt.Sprintf("%d", event.ID))
	if err != nil {
		return errs.NewCannotBeProcessedError(fmt.Sprintf("error deleting document: %v", err))
	}
	defer res.Body.Close()

	logs.Info(fmt.Sprintf("Job document deleted: %v", event))

	return nil
}

func (os *openSearchRepository) CreateOrUpdateOrganization(event *models.OrganizationDocument) error {
	// Convert organization data to JSON
	data, err := json.Marshal(event)
	if err != nil {
		return errs.NewCannotBeProcessedError("error marshalling organization data")
	}

	// Index the document
	res, err := os.es.Index("organization", bytes.NewReader(data), os.es.Index.WithDocumentID(fmt.Sprintf("%d", event.ID)))
	if err != nil {
		return errs.NewCannotBeProcessedError("error indexing document")
	}
	defer res.Body.Close()

	logs.Info(fmt.Sprintf("Organization document indexed: %v", event))

	return nil
}

func (os *openSearchRepository) DeleteOrganization(event models.OrganizationDocument) error {
	// Delete the document from OpenSearch
	res, err := os.es.Delete("organization", fmt.Sprintf("%d", event.ID))
	if err != nil {
		return errs.NewCannotBeProcessedError(fmt.Sprintf("error deleting document: %v", err))
	}
	defer res.Body.Close()

	logs.Info(fmt.Sprintf("Organization document deleted: %v", event))

	return nil
}

// func (os *openSearchRepository) CreateOrUpdateEvent(event models.CDCEvent) error {
// 	indexName := event.Payload.Source.Table
// 	if indexName != "events" && indexName != "jobs" && indexName != "organization" {
// 		return nil // Ignore other tables
// 	}

// 	// Extract the document ID
// 	docID, ok := event.Payload.After["id"].(string)
// 	if !ok {
// 		return errs.NewBadRequestError("missing or invalid document ID")
// 	}

// 	// Convert event data to JSON
// 	data, err := json.Marshal(event.Payload.After)
// 	if err != nil {
// 		return errs.NewCannotBeProcessedError("error marshalling event data")
// 	}

// 	// Index the document
// 	res, err := os.es.Index(indexName, bytes.NewReader(data), os.es.Index.WithDocumentID(docID))
// 	if err != nil {
// 		return errs.NewCannotBeProcessedError("error indexing document")
// 	}
// 	defer res.Body.Close()

// 	return nil
// }

// func (os *openSearchRepository) DeleteEvent(event models.CDCEvent) error {
// 	indexName := event.Payload.Source.Table
// 	if indexName != "events" && indexName != "jobs" && indexName != "organization" {
// 		return nil // Ignore other tables
// 	}

// 	// Extract document ID
// 	docID, ok := event.Payload.Before["id"].(string)
// 	if !ok {
// 		return errs.NewBadRequestError("missing or invalid document ID")
// 	}

// 	// Delete the document from OpenSearch
// 	res, err := os.es.Delete(indexName, docID)
// 	if err != nil {
// 		return errs.NewCannotBeProcessedError(fmt.Sprintf("error deleting document: %v", err))
// 	}
// 	defer res.Body.Close()

// 	return nil
// }
