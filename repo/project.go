package repo

import (
	"register-system-be/model"

	"gorm.io/gorm"
)

type ProjectRepo interface {
	Create(p *model.Project) error
	FindByID(id uint) (*model.Project, error)
	CountAttending(projectID uint) (uint, error)
}

type projectRepo struct {
	db *gorm.DB
}

func NewProjectRepo(db *gorm.DB) ProjectRepo {
	return &projectRepo{db: db}
}

func (r *projectRepo) Create(p *model.Project) error {
	return r.db.Create(p).Error
}

func (r *projectRepo) FindByID(id uint) (*model.Project, error) {
	var p model.Project
	if err := r.db.Preload("Affiliation").Preload("CommunityUser").First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
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
