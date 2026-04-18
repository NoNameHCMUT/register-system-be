package business

import (
	"errors"

	"register-system-be/email"
	"register-system-be/model"
	"register-system-be/repo"
)

type AdminBusiness interface {
	ListPending() ([]model.UserResponse, error)
	ListActiveUsers() ([]model.UserResponse, error)
	AcceptUser(userID uint) (*model.UserResponse, error)
	RejectUser(userID uint) error
	ListAllProjects() ([]model.ProjectResponse, error)
	ListAllApplications() ([]model.StudentProjectResponse, error)
	CreateAffiliation(req *model.AffiliationCreateRequest) (*model.AffiliationResponse, error)
	UpdateAffiliation(id uint, req *model.AffiliationUpdateRequest) (*model.AffiliationResponse, error)
	DeleteAffiliation(id uint) error
}

type adminBusiness struct {
	userRepo        repo.UserRepo
	projectRepo     repo.ProjectRepo
	applicationRepo repo.ApplicationRepo
	affiliationRepo repo.AffiliationRepo
	emailSender     email.EmailSender
}

func NewAdminBusiness(ur repo.UserRepo, pr repo.ProjectRepo, ar repo.ApplicationRepo, afr repo.AffiliationRepo, es email.EmailSender) AdminBusiness {
	return &adminBusiness{userRepo: ur, projectRepo: pr, applicationRepo: ar, affiliationRepo: afr, emailSender: es}
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
		return nil, errors.New("user already active")
	}
	if err := b.userRepo.UpdateActive(userID, true); err != nil {
		return nil, err
	}
	user.IsActive = true
	go func() { _ = email.NotifyAccountStatus(b.emailSender, user.Email, true) }()
	resp := model.ToUserResponse(user)
	return &resp, nil
}

func (b *adminBusiness) RejectUser(userID uint) error {
	user, err := b.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if err := b.userRepo.UpdateActive(userID, false); err != nil {
		return err
	}
	if err := b.userRepo.RemoveUser(userID); err != nil {
		return err
	}
	go func() { _ = email.NotifyAccountStatus(b.emailSender, user.Email, false) }()
	return nil
}

func (b *adminBusiness) ListAllProjects() ([]model.ProjectResponse, error) {
	projects, err := b.projectRepo.FindAll()
	if err != nil {
		return nil, err
	}
	res := make([]model.ProjectResponse, 0, len(projects))
	for i := range projects {
		cnt, _ := b.projectRepo.CountAttending(projects[i].ID)
		projects[i].NumAttending = cnt
		r := toProjectResponse(&projects[i])
		res = append(res, r)
	}
	return res, nil
}

func (b *adminBusiness) ListAllApplications() ([]model.StudentProjectResponse, error) {
	apps, err := b.applicationRepo.FindAll()
	if err != nil {
		return nil, err
	}
	res := make([]model.StudentProjectResponse, 0, len(apps))
	for i := range apps {
		r := toApplicationResponse(&apps[i])
		res = append(res, r)
	}
	return res, nil
}

func (b *adminBusiness) CreateAffiliation(req *model.AffiliationCreateRequest) (*model.AffiliationResponse, error) {
	aff := &model.Affiliation{
		StdName:     req.StdName,
		Description: req.Description,
	}
	if err := b.affiliationRepo.Create(aff); err != nil {
		return nil, err
	}
	return &model.AffiliationResponse{
		ID:          aff.ID,
		StdName:     aff.StdName,
		Description: aff.Description,
	}, nil
}

func (b *adminBusiness) UpdateAffiliation(id uint, req *model.AffiliationUpdateRequest) (*model.AffiliationResponse, error) {
	aff, err := b.affiliationRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("affiliation not found")
	}
	if req.StdName != nil {
		aff.StdName = *req.StdName
	}
	if req.Description != nil {
		aff.Description = *req.Description
	}
	if err := b.affiliationRepo.Update(aff); err != nil {
		return nil, err
	}
	return &model.AffiliationResponse{
		ID:          aff.ID,
		StdName:     aff.StdName,
		Description: aff.Description,
	}, nil
}

func (b *adminBusiness) DeleteAffiliation(id uint) error {
	if _, err := b.affiliationRepo.FindByID(id); err != nil {
		return errors.New("affiliation not found")
	}
	return b.affiliationRepo.Delete(id)
}
