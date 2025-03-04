package repository

import "github.com/DAF-Bridge/cdc-service/internal/models"

type OpenSearchRepository interface {
	CreateOrUpdateEvent(event *models.EventDocument) error
	DeleteEvent(event models.EventDocument) error

	CreateOrUpdateJob(event *models.JobDocument) error
	DeleteJob(event models.JobDocument) error

	CreateOrUpdateOrganization(event *models.OrganizationDocument) error
	DeleteOrganization(event models.OrganizationDocument) error
}
