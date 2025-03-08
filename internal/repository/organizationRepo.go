package repository

import (
	"github.com/DAF-Bridge/cdc-service/internal/models"
	"gorm.io/gorm"
)

type organizationRepository struct {
	db *gorm.DB
}

// Constructor
func NewOrganizationRepository(db *gorm.DB) OrganizationRepository {
	return organizationRepository{db: db}
}

func (r organizationRepository) GetAllIndustries() ([]models.Industry, error) {
	var industries []models.Industry
	err := r.db.Find(&industries).Error
	if err != nil {
		return nil, err
	}
	return industries, nil
}

func (r organizationRepository) FindIndustryByIds(industryIDs []uint) ([]models.Industry, error) {
	var industries []models.Industry
	err := r.db.Find(&industries, industryIDs).Error
	if err != nil {
		return nil, err
	}
	return industries, nil
}

func (r organizationRepository) GetByOrgID(id uint) (*models.Organization, error) {
	org := &models.Organization{}
	if err := r.db.
		Preload("OrganizationContacts").
		Preload("Industries").
		Where("id = ? ", id).
		First(org).Error; err != nil {
		return nil, err
	}
	return org, nil
}

func (r organizationRepository) GetAllOrganizations() ([]models.Organization, error) {
	var orgs []models.Organization
	err := r.db.
		Preload("OrganizationContacts").
		Preload("Industries").
		Find(&orgs).Error
	if err != nil {
		return nil, err
	}
	return orgs, nil
}

// --------------------------------------------------------------------------
// OrganizationContact Repository
// --------------------------------------------------------------------------

type organizationContactRepository struct {
	db *gorm.DB
}

func NewOrganizationContactRepository(db *gorm.DB) OrganizationContactRepository {
	return organizationContactRepository{db: db}
}

func (r organizationContactRepository) Create(orgID uint, contact *models.OrganizationContact) error {
	tx := r.db.Begin()

	contact.OrganizationID = orgID
	if err := tx.Create(contact).Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r organizationContactRepository) GetByID(orgID uint, id uint) (*models.OrganizationContact, error) {
	contact := &models.OrganizationContact{}
	if err := r.db.
		Where("organization_id = ? AND id = ?", orgID, id).
		First(contact).Error; err != nil {

		return nil, err
	}
	return contact, nil
}

func (r organizationContactRepository) GetAllByOrgID(orgID uint) ([]models.OrganizationContact, error) {
	var contancts []models.OrganizationContact
	if err := r.db.
		Where("organization_id = ?", orgID).
		Find(&contancts).Error; err != nil {

		return nil, err
	}
	return contancts, nil
}
