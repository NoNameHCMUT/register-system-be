package repo

import (
	"register-system-be/model"

	"gorm.io/gorm"
)

type ProjectRepo interface {
	GetByAffiliationID(affiliationID uint) ([]model.Project, error)
	FindByID(id uint) (*model.Project, error)
	Create(p *model.Project) error
	Update(p *model.Project) error
	CountAttending(projectID uint) (uint, error)
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

func (r *projectRepo) FindByID(id uint) (*model.Project, error) {
	var p model.Project
	if err := r.db.
		Preload("Affiliation").
		Preload("CommunityUser").
		First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *projectRepo) Create(p *model.Project) error {
	return r.db.Create(p).Error
}

func (r *projectRepo) Update(p *model.Project) error {
	return r.db.Save(p).Error
}

func (r *projectRepo) CountAttending(projectID uint) (uint, error) {
	var cnt int64
	if err := r.db.Model(&model.StudentProject{}).Where("project_id = ?", projectID).Count(&cnt).Error; err != nil {
		return 0, err
	}
	if cnt < 0 {
		cnt = 0
	}
	return uint(cnt), nil
}
