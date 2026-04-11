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

func (h *ProjectHandler) Create(c *gin.Context) {
	req, ok := Parse[model.CreateProjectRequest](c)
	if !ok {
		return
	}

	res, err := h.biz.CreateProject(GetUserID(c), req)
	if err != nil {
		Error(c, http.StatusBadRequest, err.Error())
		return
	}

	Created(c, res)
}
