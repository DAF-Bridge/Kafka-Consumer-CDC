package repository

import "github.com/DAF-Bridge/cdc-service/internal/models"

type EventRepository interface {
	GetAll() ([]models.Event, error)
	GetByID(eventID uint) (*models.Event, error)
	GetByIDwithOrgID(eventID, orgID uint) (*models.Event, error)
	FindCategoryByIds(catIDs []uint) ([]models.Category, error)
	GetAllCategories() ([]models.Category, error)
}
