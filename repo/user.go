package repo

import (
	"register-system-be/model"

	"gorm.io/gorm"
)

type UserRepo interface {
	Create(user *model.User) error
	FindByUsername(username string) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	FindByID(id uint) (*model.User, error)
	FindPending() ([]model.User, error)
	FindActive() ([]model.User, error)
	UpdateActive(userID uint, active bool) error
	UpdateRefreshToken(userID uint, token string) error
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return &userRepo{db: db}
}

func (r *userRepo) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepo) FindByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) FindByUsername(username string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) FindByID(id uint) (*model.User, error) {
	var user model.User
	if err := r.db.Preload("Affiliation").First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) FindPending() ([]model.User, error) {
	var users []model.User
	if err := r.db.Preload("Affiliation").Where("is_active = ?", false).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepo) FindActive() ([]model.User, error) {
	var users []model.User
	if err := r.db.Preload("Affiliation").Where("is_active = ?", true).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepo) UpdateActive(userID uint, active bool) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).Update("is_active", active).Error
}

func (r *userRepo) UpdateRefreshToken(userID uint, token string) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).Update("refresh_token", token).Error
}
