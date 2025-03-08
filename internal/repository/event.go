package repository

import (
	"github.com/DAF-Bridge/cdc-service/internal/models"
	"gorm.io/gorm"
)

type eventRepository struct {
	db *gorm.DB
}

// Constructor EventRepository
func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{db: db}
}

func (r eventRepository) GetAll() ([]models.Event, error) {
	var events []models.Event
	err := r.db.
		Preload("ContactChannels").
		Preload("Categories").
		Preload("Organization").
		Find(&events).Error
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (r eventRepository) GetAllByOrgID(orgID uint) ([]models.Event, error) {
	var events []models.Event

	err := r.db.
		Preload("ContactChannels").
		Preload("Categories").
		Preload("Organization").
		Where("organization_id = ?", orgID).
		Find(&events).Error
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (r eventRepository) GetByID(eventID uint) (*models.Event, error) {
	event := models.Event{}

	if err := r.db.
		Preload("Organization").
		Preload("Categories").
		Preload("ContactChannels").
		Where("id = ?", eventID).
		First(&event).Error; err != nil {

		return nil, err
	}

	return &event, nil
}

func (r eventRepository) GetByIDwithOrgID(orgID uint, eventID uint) (*models.Event, error) {
	event := models.Event{}

	err := r.db.
		Preload("Organization").
		Preload("Categories").
		Preload("ContactChannels").
		Where("organization_id = ? AND id = ?", orgID, eventID).
		First(&event).Error

	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (r eventRepository) GetAllCategories() ([]models.Category, error) {
	var categories []models.Category

	err := r.db.Find(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (r eventRepository) FindCategoryByIds(catIDs []uint) ([]models.Category, error) {
	var categories []models.Category

	err := r.db.Find(&categories, catIDs).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}
