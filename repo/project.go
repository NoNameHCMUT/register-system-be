package repo

import (
	"register-system-be/model"

	"gorm.io/gorm"
)

type ProjectRepo interface {
	Create(p *model.Project) error
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
