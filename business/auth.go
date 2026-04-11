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
	Register(req *model.RegisterRequest) (*model.AuthResponse, error)
	Login(req *model.LoginRequest) (*model.AuthResponse, error)
	Refresh(req *model.RefreshRequest) (*model.AuthResponse, error)
	GetCurrentUser(id uint) (*model.UserResponse, error)
}

type authBusiness struct {
	userRepo  repo.UserRepo
	tokenRepo repo.TokenRepo
	cfg       *config.Config
}

func NewAuthBusiness(ur repo.UserRepo, tr repo.TokenRepo, cfg *config.Config) AuthBusiness {
	return &authBusiness{userRepo: ur, tokenRepo: tr, cfg: cfg}
}

func (b *authBusiness) Register(req *model.RegisterRequest) (*model.AuthResponse, error) {
	if _, err := b.userRepo.FindByEmail(req.Email); err == nil {
		return nil, errors.New("email already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email:     req.Email,
		Password:  string(hash),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      model.RoleUser,
	}

	if err := b.userRepo.Create(user); err != nil {
		return nil, err
	}

	return b.generateTokens(user)
}

func (b *authBusiness) Login(req *model.LoginRequest) (*model.AuthResponse, error) {
	user, err := b.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return b.generateTokens(user)
}

func (b *authBusiness) Refresh(req *model.RefreshRequest) (*model.AuthResponse, error) {
	stored, err := b.tokenRepo.FindByToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if stored.Revoked {
		return nil, errors.New("refresh token revoked")
	}

	if time.Now().After(stored.ExpiresAt) {
		return nil, errors.New("refresh token expired")
	}

	user, err := b.userRepo.FindByID(stored.UserID)
	if err != nil {
		return nil, err
	}

	if err := b.tokenRepo.Revoke(req.RefreshToken); err != nil {
		return nil, err
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
	accessToken, err := b.createToken(user, b.cfg.JWTAccessExpiry)
	if err != nil {
		return nil, err
	}

	refreshToken, err := b.createToken(user, b.cfg.JWTRefreshExpiry)
	if err != nil {
		return nil, err
	}

	rt := &model.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(b.cfg.JWTRefreshExpiry),
	}
	if err := b.tokenRepo.Create(rt); err != nil {
		return nil, err
	}

	return &model.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         model.ToUserResponse(user),
	}, nil
}

func (b *authBusiness) createToken(user *model.User, expiry time.Duration) (string, error) {
	claims := model.Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(b.cfg.JWTSecret))
}
