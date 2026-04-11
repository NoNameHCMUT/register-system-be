package repo

import (
	"register-system-be/model"

	"gorm.io/gorm"
)

type ProjectRepo interface {
	GetByAffiliationID(affiliationID uint) ([]model.Project, error)
	GetByID(id uint) (*model.Project, error)
}

type projectRepo struct {
	db *gorm.DB
}

func NewProjectRepo(db *gorm.DB) ProjectRepo {
	return &projectRepo{db: db}
}

// GetByAffiliationID lấy tất cả project của một affiliation
func (r *projectRepo) GetByAffiliationID(affiliationID uint) ([]model.Project, error) {
	var projects []model.Project
	if err := r.db.
		Preload("Affiliation").
		Preload("CommunityUser").
		Where("affiliation_id = ?", affiliationID).
		Order("created_at DESC").
		Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

// GetByID lấy project theo ID
func (r *projectRepo) GetByID(id uint) (*model.Project, error) {
	var project model.Project
	if err := r.db.
		Preload("Affiliation").
		Preload("CommunityUser").
		First(&project, id).Error; err != nil {
		return nil, err
	}
	return &project, nil
}
