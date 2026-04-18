package business

import (
	"errors"
	"time"

	"register-system-be/email"
	"register-system-be/model"
	"register-system-be/repo"
)

type ApplicationBusiness interface {
	Apply(studentID uint, projectID uint) (*model.StudentProjectResponse, error)
	ListByStudent(studentID uint) ([]model.StudentProjectResponse, error)
	ListByProjectForSchool(schoolUserID uint, projectID uint) ([]model.StudentProjectResponse, error)
	ListByProjectForCommunity(communityUserID uint, projectID uint) ([]model.StudentProjectResponse, error)
	SchoolAction(schoolUserID uint, req *model.ApplicationActionRequest) (int, error)
	CommunityAction(communityUserID uint, req *model.ApplicationActionRequest) (int, error)
}

type applicationBusiness struct {
	applicationRepo repo.ApplicationRepo
	projectRepo     repo.ProjectRepo
	userRepo        repo.UserRepo
	emailSender     email.EmailSender
}

func NewApplicationBusiness(ar repo.ApplicationRepo, pr repo.ProjectRepo, ur repo.UserRepo, es email.EmailSender) ApplicationBusiness {
	return &applicationBusiness{applicationRepo: ar, projectRepo: pr, userRepo: ur, emailSender: es}
}

func (b *applicationBusiness) Apply(studentID uint, projectID uint) (*model.StudentProjectResponse, error) {
	user, err := b.userRepo.FindByID(studentID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if user.Role != model.RoleStudent {
		return nil, errors.New("only students can apply")
	}

	project, err := b.projectRepo.FindByID(projectID)
	if err != nil {
		return nil, errors.New("project not found")
	}
	if project.DateApproved == nil {
		return nil, errors.New("project not yet approved by school")
	}
	if project.AffiliationID != user.AffiliationID {
		return nil, errors.New("project not available for your affiliation")
	}

	now := time.Now()
	if now.Before(project.FormStartDay) || now.After(project.FormEndDay) {
		return nil, errors.New("not within form registration period")
	}

	cnt, err := b.projectRepo.CountAttending(projectID)
	if err != nil {
		return nil, errors.New("failed to check project capacity")
	}
	if cnt >= project.NumMax {
		return nil, errors.New("project has reached maximum number of participants")
	}

	existing, _ := b.applicationRepo.FindByStudentAndProject(studentID, projectID)
	if existing != nil {
		return nil, errors.New("already applied to this project")
	}

	sp := &model.StudentProject{
		UserID:    studentID,
		ProjectID: projectID,
		Status:    model.StatusSchoolPending,
	}

	if err := b.applicationRepo.Create(sp); err != nil {
		return nil, err
	}

	created, err := b.applicationRepo.FindByID(sp.ID)
	if err != nil {
		return nil, err
	}

	resp := toApplicationResponse(created)
	return &resp, nil
}

func (b *applicationBusiness) ListByStudent(studentID uint) ([]model.StudentProjectResponse, error) {
	user, err := b.userRepo.FindByID(studentID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if user.Role != model.RoleStudent {
		return nil, errors.New("permission denied")
	}

	apps, err := b.applicationRepo.FindByStudentID(studentID)
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

func (b *applicationBusiness) ListByProjectForSchool(schoolUserID uint, projectID uint) ([]model.StudentProjectResponse, error) {
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

	apps, err := b.applicationRepo.FindByProjectIDAndStatus(projectID, model.StatusSchoolPending)
	if err != nil {
		return nil, err
	}

	allApps, _ := b.applicationRepo.FindByProjectID(projectID)
	approvedApps := make([]model.StudentProject, 0)
	for _, a := range allApps {
		if a.Status != model.StatusSchoolPending {
			approvedApps = append(approvedApps, a)
		}
	}
	apps = append(apps, approvedApps...)

	res := make([]model.StudentProjectResponse, 0, len(apps))
	for i := range apps {
		r := toApplicationResponse(&apps[i])
		res = append(res, r)
	}
	return res, nil
}

func (b *applicationBusiness) ListByProjectForCommunity(communityUserID uint, projectID uint) ([]model.StudentProjectResponse, error) {
	user, err := b.userRepo.FindByID(communityUserID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if user.Role != model.RoleCommunity {
		return nil, errors.New("permission denied")
	}

	project, err := b.projectRepo.FindByID(projectID)
	if err != nil {
		return nil, errors.New("project not found")
	}
	if project.CommunityUserID != communityUserID {
		return nil, errors.New("not your project")
	}

	apps, err := b.applicationRepo.FindByProjectID(projectID)
	if err != nil {
		return nil, err
	}

	visible := make([]model.StudentProject, 0)
	for _, a := range apps {
		if a.Status == model.StatusCommunityPending || a.Status == model.StatusCommunityReject || a.Status == model.StatusApproved {
			visible = append(visible, a)
		}
	}

	res := make([]model.StudentProjectResponse, 0, len(visible))
	for i := range visible {
		r := toApplicationResponse(&visible[i])
		res = append(res, r)
	}
	return res, nil
}

func (b *applicationBusiness) SchoolAction(schoolUserID uint, req *model.ApplicationActionRequest) (int, error) {
	user, err := b.userRepo.FindByID(schoolUserID)
	if err != nil {
		return 0, errors.New("user not found")
	}
	if user.Role != model.RoleSchool {
		return 0, errors.New("permission denied")
	}

	var targetStatus model.ApplicationStatus
	var expectedCurrent model.ApplicationStatus

	if req.Action == "approve" {
		targetStatus = model.StatusCommunityPending
		expectedCurrent = model.StatusSchoolPending
	} else {
		targetStatus = model.StatusSchoolReject
		expectedCurrent = model.StatusSchoolPending
	}

	updated := 0
	for _, id := range req.ApplicationIDs {
		app, err := b.applicationRepo.FindByID(id)
		if err != nil {
			continue
		}
		if app.Status != expectedCurrent {
			continue
		}
		project, err := b.projectRepo.FindByID(app.ProjectID)
		if err != nil {
			continue
		}
		if project.AffiliationID != user.AffiliationID {
			continue
		}
		if err := b.applicationRepo.UpdateStatus(id, targetStatus); err != nil {
			continue
		}
		updated++

		appUser, err := b.userRepo.FindByID(app.UserID)
		if err == nil {
			go func() {
				_ = email.NotifyApplicationStatus(b.emailSender, appUser.Email, string(targetStatus), project.Name)
			}()
		}
	}

	return updated, nil
}

func (b *applicationBusiness) CommunityAction(communityUserID uint, req *model.ApplicationActionRequest) (int, error) {
	user, err := b.userRepo.FindByID(communityUserID)
	if err != nil {
		return 0, errors.New("user not found")
	}
	if user.Role != model.RoleCommunity {
		return 0, errors.New("permission denied")
	}

	var targetStatus model.ApplicationStatus

	if req.Action == "approve" {
		targetStatus = model.StatusApproved
	} else {
		targetStatus = model.StatusCommunityReject
	}

	updated := 0
	for _, id := range req.ApplicationIDs {
		app, err := b.applicationRepo.FindByID(id)
		if err != nil {
			continue
		}
		if app.Status != model.StatusCommunityPending {
			continue
		}
		project, err := b.projectRepo.FindByID(app.ProjectID)
		if err != nil {
			continue
		}
		if project.CommunityUserID != communityUserID {
			continue
		}
		if err := b.applicationRepo.UpdateStatus(id, targetStatus); err != nil {
			continue
		}
		updated++

		appUser, err := b.userRepo.FindByID(app.UserID)
		if err == nil {
			go func() {
				_ = email.NotifyApplicationStatus(b.emailSender, appUser.Email, string(targetStatus), project.Name)
			}()
		}
	}

	return updated, nil
}
