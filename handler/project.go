package handler

import (
	"net/http"

	"register-system-be/business"
	"register-system-be/repo"
	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	biz      business.ProjectBusiness
	userRepo repo.UserRepo
}

func NewProjectHandler(biz business.ProjectBusiness, userRepo repo.UserRepo) *ProjectHandler {
	return &ProjectHandler{
		biz:      biz,
		userRepo: userRepo,
	}
}

// @Summary      Get projects for student
// @Description  Get all projects available for current student's affiliation
// @Tags         Project
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} model.ProjectResponse
// @Failure      401 {object} map[string]string "Unauthorized"
// @Failure      403 {object} map[string]string "Only students can access"
// @Failure      500 {object} map[string]string "Server error"
// @Router       /projects/my-list [get]
func (h *ProjectHandler) GetMyProjects(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	projects, err := h.biz.GetStudentProjects(userID.(uint), h.userRepo)
	if err != nil {
		if err.Error() == "only students can access this endpoint" {
			Error(c, http.StatusForbidden, err.Error())
		} else if err.Error() == "user not found" {
			Error(c, http.StatusNotFound, err.Error())
		} else {
			Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	Success(c, projects)
}
