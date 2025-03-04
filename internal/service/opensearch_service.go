package service

import (
	"github.com/DAF-Bridge/cdc-service/errs"
	"github.com/DAF-Bridge/cdc-service/internal/models"
	"github.com/DAF-Bridge/cdc-service/internal/repository"
	"github.com/DAF-Bridge/cdc-service/pkg/logs"
)

type OpenSearchService struct {
	repo repository.OpenSearchRepository
}

func NewOpenSearchService(repo repository.OpenSearchRepository) *OpenSearchService {
	return &OpenSearchService{
		repo: repo,
	}
}

func (s *OpenSearchService) ProcessEvent(event models.CDCEvent) error {
	var document interface{}
	var err error

	switch event.Payload.Source.Table {
	case "events":
		logs.Info("Processing event")
		document, err = s.convertToEventDocument(event)
	case "org_open_jobs":
		logs.Info("Processing jobs")
		document, err = s.convertToJobDocument(event)
	case "organization":
		logs.Info("Processing organization")
		document, err = s.convertToOrganizationDocument(event)
	default:
		return nil // Ignore other tables
	}

	if err != nil {
		return errs.NewCannotBeProcessedError("error converting event to document")
	}

	switch doc := document.(type) {
	case *models.EventDocument:
		return s.repo.CreateOrUpdateEvent(doc)
	case *models.JobDocument:
		return s.repo.CreateOrUpdateJob(doc)
	case *models.OrganizationDocument:
		return s.repo.CreateOrUpdateOrganization(doc)
	default:
		return errs.NewCannotBeProcessedError("unknown document type")
	}
}

func (s *OpenSearchService) convertToEventDocument(event models.CDCEvent) (*models.EventDocument, error) {
	eventDoc := &models.EventDocument{
		ID:                 event.Payload.After["id"].(uint),
		Name:               event.Payload.After["name"].(string),
		Content:            event.Payload.After["content"].(string),
		Latitude:           event.Payload.After["latitude"].(float64),
		Longitude:          event.Payload.After["longitude"].(float64),
		StartDate:          event.Payload.After["start_date"].(string),
		EndDate:            event.Payload.After["end_date"].(string),
		StartTime:          event.Payload.After["start_time"].(string),
		EndTime:            event.Payload.After["end_time"].(string),
		LocationName:       event.Payload.After["location_name"].(string),
		Province:           event.Payload.After["province"].(string),
		Country:            event.Payload.After["country"].(string),
		LocationType:       event.Payload.After["location_type"].(string),
		Organization:       event.Payload.After["organization"].(string),
		OrganizationPicUrl: event.Payload.After["org_pic_url"].(string),
		Categories:         event.Payload.After["categories"].([]string),
		Audience:           event.Payload.After["audience"].(string),
		Price:              event.Payload.After["price"].(string),
	}

	return eventDoc, nil
}

func (s *OpenSearchService) convertToJobDocument(event models.CDCEvent) (*models.JobDocument, error) {
	jobDoc := &models.JobDocument{
		ID:           event.Payload.After["id"].(uint),
		Title:        event.Payload.After["title"].(string),
		Description:  event.Payload.After["description"].(string),
		Location:     event.Payload.After["location"].(string),
		Workplace:    event.Payload.After["workplace"].(string),
		WorkType:     event.Payload.After["work_type"].(string),
		CareerStage:  event.Payload.After["career_stage"].(string),
		Salary:       event.Payload.After["salary"].(float64),
		Categories:   event.Payload.After["categories"].([]string),
		Organization: event.Payload.After["organization"].(string),
		OrgPicUrl:    event.Payload.After["org_pic_url"].(string),
	}

	return jobDoc, nil
}

func (s *OpenSearchService) convertToOrganizationDocument(event models.CDCEvent) (*models.OrganizationDocument, error) {
	orgDoc := &models.OrganizationDocument{
		ID:          event.Payload.After["id"].(uint),
		Name:        event.Payload.After["org_name"].(string),
		PicUrl:      event.Payload.After["pic_url"].(string),
		Description: event.Payload.After["description"].(string),
		Location:    event.Payload.After["location"].(string),
		Email:       event.Payload.After["email"].(string),
		Phone:       event.Payload.After["phone"].(string),
	}

	return orgDoc, nil
}
