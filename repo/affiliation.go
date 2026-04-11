package repo

import (
	"register-system-be/model"

	"gorm.io/gorm"
)

type AffiliationRepo interface {
	FindByID(id uint) (*model.Affiliation, error)
	FindAll() ([]model.Affiliation, error)
}

type affiliationRepo struct {
	db *gorm.DB
}

func NewAffiliationRepo(db *gorm.DB) AffiliationRepo {
	return &affiliationRepo{db: db}
}

func (r *affiliationRepo) FindByID(id uint) (*model.Affiliation, error) {
	var aff model.Affiliation
	if err := r.db.First(&aff, id).Error; err != nil {
		return nil, err
	}
	return &aff, nil
}

func (r *affiliationRepo) FindAll() ([]model.Affiliation, error) {
	var list []model.Affiliation
	if err := r.db.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
