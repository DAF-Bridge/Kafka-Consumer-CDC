package repository

import (
	"github.com/DAF-Bridge/cdc-service/internal/models"
	"gorm.io/gorm"
)

type jobRepository struct {
	db *gorm.DB
}

func NewJobRepository(db *gorm.DB) jobRepository {
	return jobRepository{db: db}
}

func (r jobRepository) FindPreqByJobID(jobID uint) ([]models.Prerequisite, error) {
	var pre []models.Prerequisite
	if err := r.db.Where("job_id = ?", jobID).Find(&pre).Error; err != nil {
		return nil, err
	}

	return pre, nil
}

func (r jobRepository) FindCategoryByIds(catIDs []uint) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Find(&categories, catIDs).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r jobRepository) GetAllJobs() ([]models.OrgOpenJob, error) {
	var orgs []models.OrgOpenJob
	err := r.db.
		Preload("Organization").
		Preload("Prerequisites").
		Preload("Categories").
		Find(&orgs).Error
	if err != nil {
		return nil, err
	}

	return orgs, nil
}

func (r jobRepository) GetAllJobsByOrgID(OrgId uint) ([]models.OrgOpenJob, error) {
	var orgs []models.OrgOpenJob
	if err := r.db.
		Preload("Organization").
		Preload("Prerequisites").
		Preload("Categories").
		Where("organization_id = ?", OrgId).
		Find(&orgs).Error; err != nil {
		return nil, err
	}

	return orgs, nil
}

func (r jobRepository) GetJobByID(jobID uint) (*models.OrgOpenJob, error) {
	job := &models.OrgOpenJob{}

	if err := r.db.
		Preload("Organization").
		Preload("Prerequisites").
		Preload("Categories").
		Where("id = ?", jobID).
		First(&job).Error; err != nil {
		return nil, err
	}

	return job, nil
}

func (r jobRepository) GetJobByIDwithOrgID(orgID uint, jobID uint) (*models.OrgOpenJob, error) {
	job := &models.OrgOpenJob{}

	if err := r.db.
		Preload("Organization").
		Preload("Prerequisites").
		Preload("Categories").
		Where("organization_id = ? AND id = ?", orgID, jobID).
		First(&job).Error; err != nil {
		return nil, err
	}

	return job, nil
}
