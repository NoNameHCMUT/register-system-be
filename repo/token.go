package repo

import (
	"register-system-be/model"

	"gorm.io/gorm"
)

type TokenRepo interface {
	Create(token *model.RefreshToken) error
	FindByToken(tokenStr string) (*model.RefreshToken, error)
	Revoke(tokenStr string) error
	RevokeAllByUserID(userID uint) error
}

type tokenRepo struct {
	db *gorm.DB
}

func NewTokenRepo(db *gorm.DB) TokenRepo {
	return &tokenRepo{db: db}
}

func (r *tokenRepo) Create(token *model.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *tokenRepo) FindByToken(tokenStr string) (*model.RefreshToken, error) {
	var token model.RefreshToken
	if err := r.db.Where("token = ?", tokenStr).First(&token).Error; err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *tokenRepo) Revoke(tokenStr string) error {
	return r.db.Model(&model.RefreshToken{}).Where("token = ?", tokenStr).Update("revoked", true).Error
}

func (r *tokenRepo) RevokeAllByUserID(userID uint) error {
	return r.db.Model(&model.RefreshToken{}).Where("user_id = ?", userID).Update("revoked", true).Error
}
