package repository

import "github.com/DAF-Bridge/cdc-service/internal/models"

type JobRepository interface {
	FindPreqByJobID(jobID uint) ([]models.Prerequisite, error)
	FindCategoryByIds(catIDs []uint) ([]models.Category, error)
	GetJobByID(jobID uint) (*models.OrgOpenJob, error)
	GetJobByIDwithOrgID(orgID uint, jobID uint) (*models.OrgOpenJob, error)
	GetAllJobs() ([]models.OrgOpenJob, error)
	GetAllJobsByOrgID(OrgId uint) ([]models.OrgOpenJob, error)
}

type PrerequisiteRepository interface {
	CreatePrerequisite(jobID uint, prerequisite *models.Prerequisite) error
	GetPrerequisiteByID(prerequisiteID uint) (*models.Prerequisite, error)
	GetAllPrerequisites() ([]models.Prerequisite, error)
	GetAllPrerequisitesBelongToJobs(jobID uint) ([]models.Prerequisite, error)
	UpdatePrerequisite(prerequisite *models.Prerequisite) (*models.Prerequisite, error)
	DeletePrerequisite(prerequisiteID uint) error
}
