package business

import (
	"errors"
	"register-system-be/model"
	"register-system-be/repo"
)

type ProjectBusiness interface {
	GetStudentProjects(userID uint, userRepo repo.UserRepo) ([]model.ProjectResponse, error)
}

type projectBusiness struct {
	projectRepo repo.ProjectRepo
}

func NewProjectBusiness(projectRepo repo.ProjectRepo) ProjectBusiness {
	return &projectBusiness{projectRepo: projectRepo}
}

// GetStudentProjects lấy các project của student dựa trên affiliation
func (b *projectBusiness) GetStudentProjects(userID uint, userRepo repo.UserRepo) ([]model.ProjectResponse, error) {
	// 1. Lấy user info để biết affiliation
	user, err := userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 2. Kiểm tra user là student
	if user.Role != model.RoleStudent {
		return nil, errors.New("only students can access this endpoint")
	}

	// 3. Lấy tất cả project của affiliation này
	projects, err := b.projectRepo.GetByAffiliationID(user.AffiliationID)
	if err != nil {
		return nil, errors.New("failed to fetch projects")
	}

	// 4. Convert sang response
	responses := make([]model.ProjectResponse, 0, len(projects))
	for _, p := range projects {
		responses = append(responses, model.ToProjectResponse(&p))
	}

	return responses, nil
}
