package business

import (
	"errors"
	"time"

	"register-system-be/model"
	"register-system-be/repo"
)

type ProjectBusiness interface {
	Create(creatorID uint, creatorRole model.Role, req *model.ProjectCreateRequest) (*model.ProjectResponse, error)
}

type projectBusiness struct {
	projectRepo     repo.ProjectRepo
	affiliationRepo repo.AffiliationRepo
	userRepo        repo.UserRepo
}

func NewProjectBusiness(pr repo.ProjectRepo, ar repo.AffiliationRepo, ur repo.UserRepo) ProjectBusiness {
	return &projectBusiness{projectRepo: pr, affiliationRepo: ar, userRepo: ur}
}

func parseRFC3339(s string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}, errors.New("invalid time format, expected RFC3339")
	}
	return t, nil
}

func toProjectResponse(p *model.Project) model.ProjectResponse {
	resp := model.ProjectResponse{
		ID:              p.ID,
		AffiliationID:   p.AffiliationID,
		CommunityUserID: p.CommunityUserID,
		Name:            p.Name,
		Description:     p.Description,
		NumMax:          p.NumMax,
		NumAttending:    p.NumAttending,
		ProjectStartDay: p.ProjectStartDay.Format(time.RFC3339),
		ProjectEndDay:   p.ProjectEndDay.Format(time.RFC3339),
		FormStartDay:    p.FormStartDay.Format(time.RFC3339),
		FormEndDay:      p.FormEndDay.Format(time.RFC3339),
		CreatedAt:       p.CreatedAt.Format(time.RFC3339),
	}

	if p.Affiliation.ID != 0 {
		resp.Affiliation = &model.AffiliationResponse{ID: p.Affiliation.ID, StdName: p.Affiliation.StdName}
	}
	if p.CommunityUser.ID != 0 {
		u := model.ToUserResponse(&p.CommunityUser)
		resp.CommunityUser = &u
	}
	return resp
}

func (b *projectBusiness) Create(creatorID uint, creatorRole model.Role, req *model.ProjectCreateRequest) (*model.ProjectResponse, error) {
	_ = creatorRole // do not trust token role; re-check from DB
	user, err := b.userRepo.FindByID(creatorID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if user.Role != model.RoleCommunity {
		return nil, errors.New("permission denied")
	}

	if _, err := b.affiliationRepo.FindByID(req.AffiliationID); err != nil {
		return nil, errors.New("affiliation not found")
	}

	ps, err := parseRFC3339(req.ProjectStartDay)
	if err != nil {
		return nil, err
	}
	pe, err := parseRFC3339(req.ProjectEndDay)
	if err != nil {
		return nil, err
	}
	fs, err := parseRFC3339(req.FormStartDay)
	if err != nil {
		return nil, err
	}
	fe, err := parseRFC3339(req.FormEndDay)
	if err != nil {
		return nil, err
	}
	if !ps.Before(pe) {
		return nil, errors.New("project_start_day must be before project_end_day")
	}
	if !fs.Before(fe) {
		return nil, errors.New("form_start_day must be before form_end_day")
	}

	p := &model.Project{
		AffiliationID:   req.AffiliationID,
		CommunityUserID: creatorID,
		Name:            req.Name,
		Description:     req.Description,
		NumMax:          req.NumMax,
		ProjectStartDay: ps,
		ProjectEndDay:   pe,
		FormStartDay:    fs,
		FormEndDay:      fe,
	}

	if err := b.projectRepo.Create(p); err != nil {
		return nil, err
	}

	created, err := b.projectRepo.FindByID(p.ID)
	if err != nil {
		return nil, err
	}
	cnt, err := b.projectRepo.CountAttending(created.ID)
	if err != nil {
		return nil, err
	}
	created.NumAttending = cnt

	resp := toProjectResponse(created)
	return &resp, nil
}
