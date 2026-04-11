package business

import (
	"errors"

	"register-system-be/model"
	"register-system-be/repo"
)

type ProjectBusiness interface {
	CreateProject(currentUserID uint, req *model.CreateProjectRequest) (*model.ProjectResponse, error)
}

type projectBusiness struct {
	userRepo       repo.UserRepo
	affiliationRepo repo.AffiliationRepo
	projectRepo    repo.ProjectRepo
}

func NewProjectBusiness(ur repo.UserRepo, ar repo.AffiliationRepo, pr repo.ProjectRepo) ProjectBusiness {
	return &projectBusiness{userRepo: ur, affiliationRepo: ar, projectRepo: pr}
}

func (b *projectBusiness) CreateProject(currentUserID uint, req *model.CreateProjectRequest) (*model.ProjectResponse, error) {
	user, err := b.userRepo.FindByID(currentUserID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if user.Role != model.RoleCommunity {
		return nil, errors.New("only community can create project")
	}

	if _, err := b.affiliationRepo.FindByID(req.AffiliationID); err != nil {
		return nil, errors.New("affiliation not found")
	}

	if req.ProjectEndDay.Before(req.ProjectStartDay) {
		return nil, errors.New("project_end_day must be after project_start_day")
	}
	if req.FormEndDay.Before(req.FormStartDay) {
		return nil, errors.New("form_end_day must be after form_start_day")
	}

	p := &model.Project{
		AffiliationID:   req.AffiliationID,
		CommunityUserID: currentUserID,
		Name:            req.Name,
		Description:     req.Description,
		NumMax:          req.NumMax,
		ProjectStartDay: req.ProjectStartDay,
		ProjectEndDay:   req.ProjectEndDay,
		FormStartDay:    req.FormStartDay,
		FormEndDay:      req.FormEndDay,
	}

	if err := b.projectRepo.Create(p); err != nil {
		return nil, err
	}

	resp := model.ToProjectResponse(p)
	return &resp, nil
}
