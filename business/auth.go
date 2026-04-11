package business

import (
	"errors"
	"time"

	"register-system-be/config"
	"register-system-be/model"
	"register-system-be/repo"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthBusiness interface {
	Register(req *model.RegisterRequest) (*model.RegisterResponse, error)
	Login(req *model.LoginRequest) (*model.AuthResponse, error)
	Refresh(req *model.RefreshRequest) (*model.AuthResponse, error)
	GetCurrentUser(id uint) (*model.UserResponse, error)
}

type authBusiness struct {
	userRepo        repo.UserRepo
	affiliationRepo repo.AffiliationRepo
	cfg             *config.Config
}

func NewAuthBusiness(ur repo.UserRepo, ar repo.AffiliationRepo, cfg *config.Config) AuthBusiness {
	return &authBusiness{userRepo: ur, affiliationRepo: ar, cfg: cfg}
}

func (b *authBusiness) Register(req *model.RegisterRequest) (*model.RegisterResponse, error) {
	if _, err := b.userRepo.FindByUsername(req.Username); err == nil {
		return nil, errors.New("username already exists")
	}
	if _, err := b.userRepo.FindByEmail(req.Email); err == nil {
		return nil, errors.New("email already exists")
	}
	if _, err := b.affiliationRepo.FindByID(req.AffiliationID); err != nil {
		return nil, errors.New("affiliation not found")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username:      req.Username,
		FullName:      req.FullName,
		PasswordHash:  string(hash),
		Email:         req.Email,
		StudentID:     req.StudentID,
		AffiliationID: req.AffiliationID,
		Role:          req.Role,
		IsActive:      false,
	}

	if err := b.userRepo.Create(user); err != nil {
		return nil, err
	}

	created, _ := b.userRepo.FindByID(user.ID)

	return &model.RegisterResponse{
		User:    model.ToUserResponse(created),
		Message: "account pending admin approval",
	}, nil
}

func (b *authBusiness) Login(req *model.LoginRequest) (*model.AuthResponse, error) {
	user, err := b.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.IsActive {
		return nil, errors.New("account pending admin approval")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return b.generateTokens(user)
}

func (b *authBusiness) Refresh(req *model.RefreshRequest) (*model.AuthResponse, error) {
	claims := &model.Claims{}
	token, err := jwt.ParseWithClaims(req.RefreshToken, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(b.cfg.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	user, err := b.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.RefreshToken != req.RefreshToken {
		return nil, errors.New("refresh token revoked")
	}

	return b.generateTokens(user)
}

func (b *authBusiness) GetCurrentUser(id uint) (*model.UserResponse, error) {
	user, err := b.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	resp := model.ToUserResponse(user)
	return &resp, nil
}

func (b *authBusiness) generateTokens(user *model.User) (*model.AuthResponse, error) {
	accessToken, err := b.createToken(user, model.AccessToken, b.cfg.JWTAccessExpiry)
	if err != nil {
		return nil, err
	}

	refreshToken, err := b.createToken(user, model.RefreshToken, b.cfg.JWTRefreshExpiry)
	if err != nil {
		return nil, err
	}

	if err := b.userRepo.UpdateRefreshToken(user.ID, refreshToken); err != nil {
		return nil, err
	}

	return &model.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         model.ToUserResponse(user),
	}, nil
}

func (b *authBusiness) createToken(user *model.User, tokenType model.TokenType, expiry time.Duration) (string, error) {
	claims := model.Claims{
		UserID:    user.ID,
		Email:     user.Email,
		Role:      user.Role,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(b.cfg.JWTSecret))
}
