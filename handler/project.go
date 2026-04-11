package handler

import (
	"net/http"

	"register-system-be/business"
	"register-system-be/model"

	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	biz business.ProjectBusiness
}

func NewProjectHandler(biz business.ProjectBusiness) *ProjectHandler {
	return &ProjectHandler{biz: biz}
}

// @Summary      Create project
// @Description  Create a project (community/admin)
// @Tags         Project
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body model.ProjectCreateRequest true "Create project"
// @Success      201 {object} model.ProjectResponse
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Router       /projects [post]
func (h *ProjectHandler) Create(c *gin.Context) {
	req, ok := Parse[model.ProjectCreateRequest](c)
	if !ok {
		return
	}

	res, err := h.biz.Create(GetUserID(c), model.Role(GetRole(c)), req)
	if err != nil {
		if err.Error() == "permission denied" {
			Error(c, http.StatusForbidden, err.Error())
			return
		}
		Error(c, http.StatusBadRequest, err.Error())
		return
	}
	Created(c, res)
}
