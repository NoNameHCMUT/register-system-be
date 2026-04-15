package repo

import (
	"register-system-be/model"

	"gorm.io/gorm"
)

type AffiliationRepo interface {
	FindByID(id uint) (*model.Affiliation, error)
	FindAll() ([]model.Affiliation, error)
	Create(aff *model.Affiliation) error
	Update(aff *model.Affiliation) error
	Delete(id uint) error
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
	if err := r.db.Where("id > ?", 0).Order("id asc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *affiliationRepo) Create(aff *model.Affiliation) error {
	return r.db.Create(aff).Error
}

func (r *affiliationRepo) Update(aff *model.Affiliation) error {
	return r.db.Save(aff).Error
}

func (r *affiliationRepo) Delete(id uint) error {
	return r.db.Delete(&model.Affiliation{}, id).Error
}
