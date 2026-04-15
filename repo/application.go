package repo

import (
	"register-system-be/model"

	"gorm.io/gorm"
)

type ApplicationRepo interface {
	Create(sp *model.StudentProject) error
	FindByID(id uint) (*model.StudentProject, error)
	FindByStudentID(userID uint) ([]model.StudentProject, error)
	FindByProjectID(projectID uint) ([]model.StudentProject, error)
	FindByProjectIDAndStatus(projectID uint, status model.ApplicationStatus) ([]model.StudentProject, error)
	FindByStudentAndProject(userID uint, projectID uint) (*model.StudentProject, error)
	FindAll() ([]model.StudentProject, error)
	BatchUpdateStatus(ids []uint, status model.ApplicationStatus) error
	UpdateStatus(id uint, status model.ApplicationStatus) error
}

type applicationRepo struct {
	db *gorm.DB
}

func NewApplicationRepo(db *gorm.DB) ApplicationRepo {
	return &applicationRepo{db: db}
}

func (r *applicationRepo) Create(sp *model.StudentProject) error {
	return r.db.Create(sp).Error
}

func (r *applicationRepo) FindByID(id uint) (*model.StudentProject, error) {
	var sp model.StudentProject
	if err := r.db.Preload("User").Preload("Project").First(&sp, id).Error; err != nil {
		return nil, err
	}
	return &sp, nil
}

func (r *applicationRepo) FindByStudentID(userID uint) ([]model.StudentProject, error) {
	var apps []model.StudentProject
	if err := r.db.Preload("User").Preload("Project.Affiliation").Preload("Project.CommunityUser").
		Where("user_id = ?", userID).
		Find(&apps).Error; err != nil {
		return nil, err
	}
	return apps, nil
}

func (r *applicationRepo) FindByProjectID(projectID uint) ([]model.StudentProject, error) {
	var apps []model.StudentProject
	if err := r.db.Preload("User").Preload("Project").
		Where("project_id = ?", projectID).
		Find(&apps).Error; err != nil {
		return nil, err
	}
	return apps, nil
}

func (r *applicationRepo) FindByProjectIDAndStatus(projectID uint, status model.ApplicationStatus) ([]model.StudentProject, error) {
	var apps []model.StudentProject
	if err := r.db.Preload("User").Preload("Project").
		Where("project_id = ? AND status = ?", projectID, status).
		Find(&apps).Error; err != nil {
		return nil, err
	}
	return apps, nil
}

func (r *applicationRepo) FindByStudentAndProject(userID uint, projectID uint) (*model.StudentProject, error) {
	var sp model.StudentProject
	if err := r.db.Where("user_id = ? AND project_id = ?", userID, projectID).First(&sp).Error; err != nil {
		return nil, err
	}
	return &sp, nil
}

func (r *applicationRepo) FindAll() ([]model.StudentProject, error) {
	var apps []model.StudentProject
	if err := r.db.Preload("User").Preload("Project").Find(&apps).Error; err != nil {
		return nil, err
	}
	return apps, nil
}

func (r *applicationRepo) BatchUpdateStatus(ids []uint, status model.ApplicationStatus) error {
	return r.db.Model(&model.StudentProject{}).Where("id IN ?", ids).Update("status", status).Error
}

func (r *applicationRepo) UpdateStatus(id uint, status model.ApplicationStatus) error {
	return r.db.Model(&model.StudentProject{}).Where("id = ?", id).Update("status", status).Error
}
