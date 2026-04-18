package business

import (
	"register-system-be/model"
	"register-system-be/repo"
)

type AdminBusiness interface {
	ListPending() ([]model.UserResponse, error)
	ListActiveUsers() ([]model.UserResponse, error)
	AcceptUser(userID uint) (*model.UserResponse, error)
	RejectUser(userID uint) error
}

type adminBusiness struct {
	userRepo repo.UserRepo
}

func NewAdminBusiness(ur repo.UserRepo) AdminBusiness {
	return &adminBusiness{userRepo: ur}
}

func (b *adminBusiness) ListPending() ([]model.UserResponse, error) {
	users, err := b.userRepo.FindPending()
	if err != nil {
		return nil, err
	}
	res := make([]model.UserResponse, len(users))
	for i, u := range users {
		res[i] = model.ToUserResponse(&u)
	}
	return res, nil
}

func (b *adminBusiness) ListActiveUsers() ([]model.UserResponse, error) {
	users, err := b.userRepo.FindActive()
	if err != nil {
		return nil, err
	}
	res := make([]model.UserResponse, len(users))
	for i, u := range users {
		res[i] = model.ToUserResponse(&u)
	}
	return res, nil
}

func (b *adminBusiness) AcceptUser(userID uint) (*model.UserResponse, error) {
	user, err := b.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user.IsActive {
		return nil, err
	}
	if err := b.userRepo.UpdateActive(userID, true); err != nil {
		return nil, err
	}
	user.IsActive = true
	resp := model.ToUserResponse(user)
	return &resp, nil
}

func (b *adminBusiness) RejectUser(userID uint) error {
	return b.userRepo.UpdateActive(userID, false)
}
