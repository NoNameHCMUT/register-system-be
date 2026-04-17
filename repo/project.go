package repo

import (
	"register-system-be/model"

	"gorm.io/gorm"
)

type ProjectRepo interface {
	Create(p *model.Project) error
	FindByID(id uint) (*model.Project, error)
	FindAll() ([]model.Project, error)
	FindByAffiliation(affiliationID uint) ([]model.Project, error)
	FindByAffiliationApproved(affiliationID uint) ([]model.Project, error)
	FindByCreator(userID uint) ([]model.Project, error)
	GetByAffiliationID(affiliationID uint) ([]model.Project, error)
	Update(p *model.Project) error
	UpdateDateApproved(id uint, dateApproved interface{}) error
	UpdateBanner(id uint, bannerURL string) error
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

func (r *projectRepo) FindAll() ([]model.Project, error) {
	var projects []model.Project
	if err := r.db.Preload("Affiliation").Preload("CommunityUser").Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *projectRepo) FindByAffiliation(affiliationID uint) ([]model.Project, error) {
	var projects []model.Project
	if err := r.db.Preload("Affiliation").Preload("CommunityUser").
		Where("affiliation_id = ?", affiliationID).
		Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *projectRepo) FindByAffiliationApproved(affiliationID uint) ([]model.Project, error) {
	var projects []model.Project
	if err := r.db.Preload("Affiliation").Preload("CommunityUser").
		Where("affiliation_id = ? AND date_approved IS NOT NULL", affiliationID).
		Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *projectRepo) FindByCreator(userID uint) ([]model.Project, error) {
	var projects []model.Project
	if err := r.db.Preload("Affiliation").Preload("CommunityUser").
		Where("community_user_id = ?", userID).
		Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *projectRepo) GetByAffiliationID(affiliationID uint) ([]model.Project, error) {
	return r.FindByAffiliation(affiliationID)
}

func (r *projectRepo) Update(p *model.Project) error {
	return r.db.Save(p).Error
}

func (r *projectRepo) UpdateDateApproved(id uint, dateApproved interface{}) error {
	return r.db.Model(&model.Project{}).Where("id = ?", id).Update("date_approved", dateApproved).Error
}

func (r *projectRepo) UpdateBanner(id uint, bannerURL string) error {
	return r.db.Model(&model.Project{}).Where("id = ?", id).Update("banner_url", bannerURL).Error
}

func (r *projectRepo) CountAttending(projectID uint) (uint, error) {
	var cnt int64
	if err := r.db.Model(&model.StudentProject{}).Where("project_id = ? AND status != ?", projectID, model.StatusSchoolReject).
		Count(&cnt).Error; err != nil {
		return 0, err
	}
	if cnt < 0 {
		cnt = 0
	}
	return uint(cnt), nil
}
