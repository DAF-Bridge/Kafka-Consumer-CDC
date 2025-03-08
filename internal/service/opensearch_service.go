package service

import (
	"errors"
	"fmt"

	"github.com/DAF-Bridge/cdc-service/errs"
	"github.com/DAF-Bridge/cdc-service/internal/models"
	"github.com/DAF-Bridge/cdc-service/internal/repository"
	"github.com/DAF-Bridge/cdc-service/pkg/logs"
	"gorm.io/gorm"
)

type OpenSearchService struct {
	opnRepo   repository.OpenSearchRepository
	jobRepo   repository.JobRepository
	eventRepo repository.EventRepository
	orgRepo   repository.OrganizationRepository
}

func NewOpenSearchService(opnRepo repository.OpenSearchRepository, jobRepo repository.JobRepository, eventRepo repository.EventRepository, orgRepo repository.OrganizationRepository) *OpenSearchService {
	return &OpenSearchService{
		opnRepo:   opnRepo,
		jobRepo:   jobRepo,
		eventRepo: eventRepo,
		orgRepo:   orgRepo,
	}
}

func (s *OpenSearchService) ProcessEvent(event models.CDCEvent) error {
	var document interface{}
	var err error

	switch event.Payload.Op {
	case "d":
		// Delete the document
		switch event.Payload.Source.Table {
		case "events":
			logs.Info("Deleting event")
			document = &models.EventDocument{
				ID: uint(event.Payload.Before["id"].(float64)),
			}
		case "org_open_jobs":
			logs.Info("Deleting job")
			document = &models.JobDocument{
				ID: uint(event.Payload.Before["id"].(float64)),
			}
		case "organization":
			logs.Info("Deleting organization")
			document = &models.OrganizationDocument{
				ID: uint(event.Payload.Before["id"].(float64)),
			}
		}
	case "c", "u":
		switch event.Payload.Source.Table {
		case "events":
			logs.Info("Processing event")
			document, err = s.convertToEventDocument(event)
			if err != nil {
				logs.Error(fmt.Sprintf("Error converting event to document: %v", err))
			}
		case "org_open_jobs":
			logs.Info("Processing jobs")
			document, err = s.convertToJobDocument(event)
			if err != nil {
				logs.Error(fmt.Sprintf("Error converting job to document: %v", err))
			}
		case "organization":
			logs.Info("Processing organization")
			document, err = s.convertToOrganizationDocument(event)
			if err != nil {
				logs.Error(fmt.Sprintf("Error converting organization to document: %v", err))
			}
		default:
			return nil // Ignore other tables
		}
	}

	if err != nil {
		return errs.NewCannotBeProcessedError("error converting event to document")
	}

	switch event.Payload.Op {
	case "c", "u":
		switch doc := document.(type) {
		case *models.EventDocument:
			return s.opnRepo.CreateOrUpdateEvent(doc)
		case *models.JobDocument:
			return s.opnRepo.CreateOrUpdateJob(doc)
		case *models.OrganizationDocument:
			return s.opnRepo.CreateOrUpdateOrganization(doc)
		default:
			return errs.NewCannotBeProcessedError("unknown document type")
		}
	case "d":
		switch doc := document.(type) {
		case *models.EventDocument:
			return s.opnRepo.DeleteEvent(*doc)
		case *models.JobDocument:
			return s.opnRepo.DeleteJob(*doc)
		case *models.OrganizationDocument:
			return s.opnRepo.DeleteOrganization(*doc)
		default:
			return errs.NewCannotBeProcessedError("unknown document type")
		}
	}

	return nil
}

func (s *OpenSearchService) convertToEventDocument(event models.CDCEvent) (*models.EventDocument, error) {
	id, ok := event.Payload.After["id"].(float64)
	if !ok {
		return nil, fmt.Errorf("expected id to be a float64, got %T", event.Payload.After["id"])
	}

	eventData, err := s.getEventByID(uint(id))
	if err != nil {
		logs.Error(fmt.Sprintf("Error fetching event data: %v", err))
		return nil, err
	}

	var categories []string
	for _, category := range eventData.Categories {
		categories = append(categories, category.Name)
	}

	eventDoc := &models.EventDocument{
		ID:           uint(id),
		Name:         event.Payload.After["name"].(string),
		PicUrl:       event.Payload.After["pic_url"].(string),
		Content:      event.Payload.After["content"].(string),
		Latitude:     eventData.Latitude,
		Longitude:    eventData.Longitude,
		StartDate:    eventData.StartDate.Format("2006-01-02"),
		EndDate:      eventData.EndDate.Format("2006-01-02"),
		StartTime:    eventData.StartTime.Format("15:04:05"),
		EndTime:      eventData.EndTime.Format("15:04:05"),
		LocationName: event.Payload.After["location_name"].(string),
		Province:     event.Payload.After["province"].(string),
		Country:      event.Payload.After["country"].(string),
		LocationType: event.Payload.After["location_type"].(string),
		Organization: models.OrganizationShortDocument{
			ID:     eventData.Organization.ID,
			Name:   eventData.Organization.Name,
			PicUrl: eventData.Organization.PicUrl,
		},
		Categories: categories,
		Audience:   event.Payload.After["audience"].(string),
		Price:      event.Payload.After["price_type"].(string),
		UpdateAt:   eventData.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return eventDoc, nil
}

func (s *OpenSearchService) getEventByID(id uint) (*models.Event, error) {
	event, err := s.eventRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFoundError(fmt.Sprintf("event with ID %d not found", id))
		}

		logs.Error(fmt.Sprintf("Error fetching event: %v", err))
		return nil, err
	}

	return event, nil
}

func (s *OpenSearchService) convertToJobDocument(event models.CDCEvent) (*models.JobDocument, error) {
	id, ok := event.Payload.After["id"].(float64)
	if !ok {
		return nil, fmt.Errorf("expected id to be a float64, got %T", event.Payload.After["id"])
	}

	jobData, err := s.getJobByID(uint(id))
	if err != nil {
		logs.Error(fmt.Sprintf("Error fetching job data: %v", err))
		return nil, err
	}

	var categories []string
	for _, category := range jobData.Categories {
		categories = append(categories, category.Name)
	}

	jobDoc := &models.JobDocument{
		ID:          uint(id),
		Title:       event.Payload.After["title"].(string),
		Description: event.Payload.After["description"].(string),
		Location:    event.Payload.After["location"].(string),
		Workplace:   event.Payload.After["workplace"].(string),
		WorkType:    event.Payload.After["work_type"].(string),
		CareerStage: event.Payload.After["career_stage"].(string),
		Salary:      jobData.Salary,
		Categories:  categories,
		Organization: models.OrganizationShortDocument{
			ID:     jobData.Organization.ID,
			Name:   jobData.Organization.Name,
			PicUrl: jobData.Organization.PicUrl,
		},
	}

	return jobDoc, nil
}

func (s *OpenSearchService) getJobByID(id uint) (*models.OrgOpenJob, error) {
	job, err := s.jobRepo.GetJobByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFoundError(fmt.Sprintf("job with ID %d not found", id))
		}

		logs.Error(fmt.Sprintf("Error fetching job data: %v", err))
		return nil, err
	}

	return job, nil
}

func (s *OpenSearchService) convertToOrganizationDocument(event models.CDCEvent) (*models.OrganizationDocument, error) {
	id, ok := event.Payload.After["id"].(float64)
	if !ok {
		return nil, fmt.Errorf("expected id to be a float64, got %T", event.Payload.After["id"])
	}

	orgData, err := s.getOrganizationByID(uint(id))
	if err != nil {
		logs.Error(fmt.Sprintf("Error fetching organization data: %v", err))
		return nil, err
	}

	orgDoc := &models.OrganizationDocument{
		ID:          uint(id),
		Name:        event.Payload.After["org_name"].(string),
		PicUrl:      event.Payload.After["pic_url"].(string),
		Description: event.Payload.After["description"].(string),
		Latitude:    orgData.Latitude,
		Longitude:   orgData.Longitude,
		Email:       event.Payload.After["email"].(string),
		Phone:       event.Payload.After["phone"].(string),
		Province:    event.Payload.After["province"].(string),
		Country:     event.Payload.After["country"].(string),
		UpdateAt:    orgData.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return orgDoc, nil
}

func (s *OpenSearchService) getOrganizationByID(id uint) (*models.Organization, error) {
	org, err := s.orgRepo.GetByOrgID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFoundError(fmt.Sprintf("organization with ID %d not found", id))
		}

		logs.Error(fmt.Sprintf("Error fetching organization data: %v", err))
		return nil, err
	}

	return org, nil
}
