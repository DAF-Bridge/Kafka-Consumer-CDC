package repository

import "github.com/DAF-Bridge/cdc-service/internal/models"

type OrganizationRepository interface {
	FindIndustryByIds(industryIDs []uint) ([]models.Industry, error)
	GetAllIndustries() ([]models.Industry, error)
	GetByOrgID(id uint) (*models.Organization, error)
	//GetByOrgID(userID uuid.UUID, id uint) (*models.Organization, error)
	GetAllOrganizations() ([]models.Organization, error)
	//GetAllOrganizations(userID uuid.UUID) ([]models.Organization, error)
}

type OrganizationContactRepository interface {
	GetByID(orgID uint, id uint) (*models.OrganizationContact, error)
	GetAllByOrgID(orgID uint) ([]models.OrganizationContact, error)
}
