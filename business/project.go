package business

import (
	"errors"
	"time"

	"register-system-be/model"
	"register-system-be/repo"
)

type ProjectBusiness interface {
	Create(creatorID uint, creatorRole model.Role, req *model.ProjectCreateRequest) (*model.ProjectResponse, error)
	Update(projectID uint, actorID uint, actorRole model.Role, req *model.ProjectUpdateRequest) (*model.ProjectResponse, error)
	UpdateBanner(projectID uint, actorID uint, actorRole model.Role, bannerURL string) (*model.ProjectResponse, error)
	ListBySchool(schoolUserID uint) ([]model.ProjectResponse, error)
	ListPendingBySchool(schoolUserID uint) ([]model.ProjectResponse, error)
	ApproveProject(projectID uint, schoolUserID uint) (*model.ProjectResponse, error)
	ListByCommunity(communityUserID uint) ([]model.ProjectResponse, error)
	ListApprovedForStudent(studentUserID uint) ([]model.ProjectResponse, error)
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
		BannerURL:       p.BannerURL,
		CreatedAt:       p.CreatedAt.Format(time.RFC3339),
	}

	if p.DateApproved != nil {
		resp.DateApproved = p.DateApproved.Format(time.RFC3339)
	}

	if p.Affiliation.ID != 0 {
		resp.Affiliation = &model.AffiliationResponse{ID: p.Affiliation.ID, StdName: p.Affiliation.StdName, Description: p.Affiliation.Description}
	}
	if p.CommunityUser.ID != 0 {
		u := model.ToUserResponse(&p.CommunityUser)
		resp.CommunityUser = &u
	}
	return resp
}

func toApplicationResponse(sp *model.StudentProject) model.StudentProjectResponse {
	resp := model.StudentProjectResponse{
		ID:        sp.ID,
		UserID:    sp.UserID,
		ProjectID: sp.ProjectID,
		Status:    sp.Status,
		CreatedAt: sp.CreatedAt.Format(time.RFC3339),
	}
	if sp.User.ID != 0 {
		u := model.ToUserResponse(&sp.User)
		resp.User = &u
	}
	if sp.Project.ID != 0 {
		p := toProjectResponse(&sp.Project)
		resp.Project = &p
	}
	return resp
}

func (b *projectBusiness) Create(creatorID uint, creatorRole model.Role, req *model.ProjectCreateRequest) (*model.ProjectResponse, error) {
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

func (b *projectBusiness) Update(projectID uint, actorID uint, actorRole model.Role, req *model.ProjectUpdateRequest) (*model.ProjectResponse, error) {
	user, err := b.userRepo.FindByID(actorID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	actorRole = user.Role

	p, err := b.projectRepo.FindByID(projectID)
	if err != nil {
		return nil, err
	}

	if actorRole != model.RoleAdmin && actorID != p.CommunityUserID {
		return nil, errors.New("permission denied")
	}

	if req.Name != nil {
		p.Name = *req.Name
	}
	if req.Description != nil {
		p.Description = *req.Description
	}
	if req.NumMax != nil {
		p.NumMax = *req.NumMax
	}

	if req.ProjectStartDay != nil {
		t, err := parseRFC3339(*req.ProjectStartDay)
		if err != nil {
			return nil, err
		}
		p.ProjectStartDay = t
	}
	if req.ProjectEndDay != nil {
		t, err := parseRFC3339(*req.ProjectEndDay)
		if err != nil {
			return nil, err
		}
		p.ProjectEndDay = t
	}
	if req.FormStartDay != nil {
		t, err := parseRFC3339(*req.FormStartDay)
		if err != nil {
			return nil, err
		}
		p.FormStartDay = t
	}
	if req.FormEndDay != nil {
		t, err := parseRFC3339(*req.FormEndDay)
		if err != nil {
			return nil, err
		}
		p.FormEndDay = t
	}

	if !p.ProjectStartDay.Before(p.ProjectEndDay) {
		return nil, errors.New("project_start_day must be before project_end_day")
	}
	if !p.FormStartDay.Before(p.FormEndDay) {
		return nil, errors.New("form_start_day must be before form_end_day")
	}

	if err := b.projectRepo.Update(p); err != nil {
		return nil, err
	}

	updated, err := b.projectRepo.FindByID(projectID)
	if err != nil {
		return nil, err
	}
	cnt, err := b.projectRepo.CountAttending(updated.ID)
	if err != nil {
		return nil, err
	}
	updated.NumAttending = cnt

	resp := toProjectResponse(updated)
	return &resp, nil
}

func (b *projectBusiness) UpdateBanner(projectID uint, actorID uint, actorRole model.Role, bannerURL string) (*model.ProjectResponse, error) {
	user, err := b.userRepo.FindByID(actorID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	p, err := b.projectRepo.FindByID(projectID)
	if err != nil {
		return nil, errors.New("project not found")
	}

	if user.Role != model.RoleAdmin && actorID != p.CommunityUserID {
		return nil, errors.New("permission denied")
	}

	if err := b.projectRepo.UpdateBanner(projectID, bannerURL); err != nil {
		return nil, err
	}

	updated, err := b.projectRepo.FindByID(projectID)
	if err != nil {
		return nil, err
	}
	cnt, err := b.projectRepo.CountAttending(updated.ID)
	if err != nil {
		return nil, err
	}
	updated.NumAttending = cnt

	resp := toProjectResponse(updated)
	return &resp, nil
}

func (b *projectBusiness) ListBySchool(schoolUserID uint) ([]model.ProjectResponse, error) {
	user, err := b.userRepo.FindByID(schoolUserID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if user.Role != model.RoleSchool {
		return nil, errors.New("permission denied")
	}

	projects, err := b.projectRepo.FindByAffiliation(user.AffiliationID)
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

func (b *projectBusiness) ListPendingBySchool(schoolUserID uint) ([]model.ProjectResponse, error) {
	user, err := b.userRepo.FindByID(schoolUserID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if user.Role != model.RoleSchool {
		return nil, errors.New("permission denied")
	}

	projects, err := b.projectRepo.FindByAffiliation(user.AffiliationID)
	if err != nil {
		return nil, err
	}

	res := make([]model.ProjectResponse, 0)
	for i := range projects {
		if projects[i].DateApproved != nil {
			continue
		}
		cnt, _ := b.projectRepo.CountAttending(projects[i].ID)
		projects[i].NumAttending = cnt
		r := toProjectResponse(&projects[i])
		res = append(res, r)
	}
	return res, nil
}

func (b *projectBusiness) ApproveProject(projectID uint, schoolUserID uint) (*model.ProjectResponse, error) {
	user, err := b.userRepo.FindByID(schoolUserID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if user.Role != model.RoleSchool {
		return nil, errors.New("permission denied")
	}

	project, err := b.projectRepo.FindByID(projectID)
	if err != nil {
		return nil, errors.New("project not found")
	}
	if project.AffiliationID != user.AffiliationID {
		return nil, errors.New("project not in your affiliation")
	}
	if project.DateApproved != nil {
		return nil, errors.New("project already approved")
	}

	now := time.Now()
	if err := b.projectRepo.UpdateDateApproved(projectID, now); err != nil {
		return nil, err
	}

	updated, err := b.projectRepo.FindByID(projectID)
	if err != nil {
		return nil, err
	}
	cnt, _ := b.projectRepo.CountAttending(updated.ID)
	updated.NumAttending = cnt

	resp := toProjectResponse(updated)
	return &resp, nil
}

func (b *projectBusiness) ListByCommunity(communityUserID uint) ([]model.ProjectResponse, error) {
	user, err := b.userRepo.FindByID(communityUserID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if user.Role != model.RoleCommunity {
		return nil, errors.New("permission denied")
	}

	projects, err := b.projectRepo.FindByCreator(communityUserID)
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

func (b *projectBusiness) ListApprovedForStudent(studentUserID uint) ([]model.ProjectResponse, error) {
	user, err := b.userRepo.FindByID(studentUserID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if user.Role != model.RoleStudent {
		return nil, errors.New("permission denied")
	}

	projects, err := b.projectRepo.FindByAffiliationApproved(user.AffiliationID)
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
