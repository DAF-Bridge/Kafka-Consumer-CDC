package service

import (
	"fmt"
	"time"

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
			return s.repo.CreateOrUpdateEvent(doc)
		case *models.JobDocument:
			return s.repo.CreateOrUpdateJob(doc)
		case *models.OrganizationDocument:
			return s.repo.CreateOrUpdateOrganization(doc)
		default:
			return errs.NewCannotBeProcessedError("unknown document type")
		}
	case "d":
		switch doc := document.(type) {
		case *models.EventDocument:
			return s.repo.DeleteEvent(*doc)
		case *models.JobDocument:
			return s.repo.DeleteJob(*doc)
		case *models.OrganizationDocument:
			return s.repo.DeleteOrganization(*doc)
		default:
			return errs.NewCannotBeProcessedError("unknown document type")
		}
	}

	// switch doc := document.(type) {
	// 		case *models.EventDocument:
	// 			return s.repo.CreateOrUpdateEvent(doc)
	// 		case *models.JobDocument:
	// 			return s.repo.CreateOrUpdateJob(doc)
	// 		case *models.OrganizationDocument:
	// 			return s.repo.CreateOrUpdateOrganization(doc)
	// 		default:
	// 			return errs.NewCannotBeProcessedError("unknown document type")
	// 		}
	// 	}

	return nil
}

func (s *OpenSearchService) convertToEventDocument(event models.CDCEvent) (*models.EventDocument, error) {
	id, ok := event.Payload.After["id"].(float64)
	if !ok {
		return nil, fmt.Errorf("expected id to be a float64, got %T", event.Payload.After["id"])
	}

	// latitudeStr, ok := event.Payload.After["latitude"].(string)
	// if !ok {
	// 	return nil, fmt.Errorf("expected latitude to be a string, got %T", event.Payload.After["latitude"])
	// }
	// latitudeDecoded, err := base64.StdEncoding.DecodeString(latitudeStr)
	// if err != nil {
	// 	return nil, fmt.Errorf("error decoding latitude: %v", err)
	// }
	// latitudeValue, err := strconv.ParseFloat(string(latitudeDecoded), 64)
	// if err != nil {
	// 	return nil, fmt.Errorf("error converting latitude to float64: %v", err)
	// }

	// longitudeStr, ok := event.Payload.After["longitude"].(string)
	// if !ok {
	// 	return nil, fmt.Errorf("unexpected longitude type: %T", event.Payload.After["longitude"])
	// }
	// longitudeDecoded, err := base64.StdEncoding.DecodeString(longitudeStr)
	// if err != nil {
	// 	return nil, fmt.Errorf("error decoding latitude: %v", err)
	// }
	// longitudeValue, err := parseFloat(longitudeDecoded)
	// if err != nil {
	// 	return nil, fmt.Errorf("error converting longitude to float64: %v", err)
	// }

	// startDate, ok := event.Payload.After["start_date"].(float64)
	// if !ok {
	// 	return nil, fmt.Errorf("expected start_date to be a string, got %T", event.Payload.After["start_date"])
	// }

	// EndDate, ok := event.Payload.After["end_date"].(float64)
	// if !ok {
	// 	return nil, fmt.Errorf("expected end_date to be a string, got %T", event.Payload.After["end_date"])
	// }

	// startDateValue, err := time.Parse("2006-01-02", fmt.Sprintf("%v", startDate))
	// if err != nil {
	// 	return nil, fmt.Errorf("error parsing start_date: %v", err)
	// }

	// endDateValue, err := time.Parse("2006-01-02", fmt.Sprintf("%v", EndDate))
	// if err != nil {
	// 	return nil, fmt.Errorf("error parsing end_date: %v", err)
	// }

	// categories, ok := event.Payload.After["categories"].([]string)
	// if !ok {
	// 	return nil, fmt.Errorf("expected categories to be a []string, got %T", event.Payload.After["categories"])
	// }

	// org, ok := event.Payload.After["organization"].(string)
	// if !ok {
	// 	return nil, fmt.Errorf("expected organization to be a string, got %T", event.Payload.After["organization"])
	// }

	// orgPic, ok := event.Payload.After["org_pic_url"].(string)
	// if !ok {
	// 	return nil, fmt.Errorf("expected org_pic_url to be a string, got %T", event.Payload.After["org_pic_url"])
	// }

	// type EventDocument struct {
	// 	ID           uint                      `json:"id"`
	// 	Name         string                    `json:"name"`
	// 	PicUrl       string                    `json:"picUrl"`
	// 	Content      string                    `json:"content"`
	// 	Latitude     float64                   `json:"latitude"`
	// 	Longitude    float64                   `json:"longitude"`
	// 	StartDate    string                    `json:"startDate"`
	// 	StartTime    string                    `json:"startTime"`
	// 	EndTime      string                    `json:"endTime"`
	// 	EndDate      string                    `json:"endDate"`
	// 	LocationName string                    `json:"locationName"`
	// 	Province     string                    `json:"province"`
	// 	Country      string                    `json:"country"`
	// 	LocationType string                    `json:"locationType"`
	// 	Organization OrganizationShortDocument `json:"organization"`
	// 	Categories   []string                  `json:"categories"`
	// 	Audience     string                    `json:"audience"`
	// 	Price        string                    `json:"price"`
	// 	UpdateAt     string                    `json:"updatedAt"`
	// }

	// type OrganizationShortDocument struct {
	// 	ID     uint   `json:"id"`
	// 	Name   string `json:"name"`
	// 	PicUrl string `json:"picUrl"`
	// }

	eventDoc := &models.EventDocument{
		ID:           uint(id),
		Name:         event.Payload.After["name"].(string),
		PicUrl:       event.Payload.After["pic_url"].(string),
		Content:      event.Payload.After["content"].(string),
		Latitude:     18.7964,
		Longitude:    98.9796,
		StartDate:    time.Now().Format("2006-01-02"),
		EndDate:      time.Now().Format("2006-01-02"),
		StartTime:    time.Unix(int64(event.Payload.After["start_time"].(float64)), 0).Format(time.RFC3339),
		EndTime:      time.Unix(int64(event.Payload.After["end_time"].(float64)), 0).Format(time.RFC3339),
		LocationName: event.Payload.After["location_name"].(string),
		Province:     event.Payload.After["province"].(string),
		Country:      event.Payload.After["country"].(string),
		LocationType: event.Payload.After["location_type"].(string),
		Organization: models.OrganizationShortDocument{
			ID:     uint(1),
			Name:   "Anda",
			PicUrl: "https://www.andatech.com.au/wp-content/uploads/2020/06/Andatech-Logo-2020-1.png",
		},
		Categories: []string{"Tech", "Health"},
		Audience:   event.Payload.After["audience"].(string),
		Price:      event.Payload.After["price_type"].(string),
	}

	return eventDoc, nil
}

func (s *OpenSearchService) convertToJobDocument(event models.CDCEvent) (*models.JobDocument, error) {
	id, ok := event.Payload.After["id"].(float64)
	if !ok {
		return nil, fmt.Errorf("expected id to be a float64, got %T", event.Payload.After["id"])
	}

	jobDoc := &models.JobDocument{
		ID:           uint(id),
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
	id, ok := event.Payload.After["id"].(float64)
	if !ok {
		return nil, fmt.Errorf("expected id to be a float64, got %T", event.Payload.After["id"])
	}

	orgDoc := &models.OrganizationDocument{
		ID:          uint(id),
		Name:        event.Payload.After["org_name"].(string),
		PicUrl:      event.Payload.After["pic_url"].(string),
		Description: event.Payload.After["description"].(string),
		Location:    event.Payload.After["location"].(string),
		Email:       event.Payload.After["email"].(string),
		Phone:       event.Payload.After["phone"].(string),
	}

	return orgDoc, nil
}

// func parseFloat(value interface{}) (float64, error) {
// 	switch v := value.(type) {
// 	case float64:
// 		return v, nil
// 	case string:
// 		// First, try direct float conversion
// 		if f, err := strconv.ParseFloat(v, 64); err == nil {
// 			return f, nil
// 		}
// 		// If it fails, try base64 decoding
// 		decoded, err := base64.StdEncoding.DecodeString(v)
// 		if err == nil {
// 			return strconv.ParseFloat(string(decoded), 64)
// 		}
// 		return 0, fmt.Errorf("invalid latitude/longitude format: %v", value)
// 	default:
// 		return 0, fmt.Errorf("unexpected type for latitude/longitude: %T", value)
// 	}
// }
