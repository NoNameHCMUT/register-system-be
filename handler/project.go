package handler

import (
	"net/http"
	"strconv"

	"register-system-be/business"
	"register-system-be/model"
	"register-system-be/upload"

	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	biz       business.ProjectBusiness
	uploadDir string
}

func NewProjectHandler(biz business.ProjectBusiness, uploadDir string) *ProjectHandler {
	return &ProjectHandler{biz: biz, uploadDir: uploadDir}
}

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

func (h *ProjectHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		Error(c, http.StatusBadRequest, "invalid project id")
		return
	}

	req, ok := Parse[model.ProjectUpdateRequest](c)
	if !ok {
		return
	}

	res, err := h.biz.Update(uint(id), GetUserID(c), model.Role(GetRole(c)), req)
	if err != nil {
		switch err.Error() {
		case "permission denied":
			Error(c, http.StatusForbidden, err.Error())
		case "record not found":
			Error(c, http.StatusNotFound, "project not found")
		default:
			Error(c, http.StatusBadRequest, err.Error())
		}
		return
	}

	Success(c, res)
}

func (h *ProjectHandler) UploadBanner(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		Error(c, http.StatusBadRequest, "invalid project id")
		return
	}

	file, header, err := c.Request.FormFile("banner")
	if err != nil {
		Error(c, http.StatusBadRequest, "banner file required")
		return
	}
	defer file.Close()

	path, err := upload.Save(file, header, "banners", h.uploadDir)
	if err != nil {
		Error(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.biz.UpdateBanner(uint(id), GetUserID(c), model.Role(GetRole(c)), path)
	if err != nil {
		Error(c, http.StatusBadRequest, err.Error())
		return
	}

	Success(c, res)
}

func (h *ProjectHandler) ListBySchool(c *gin.Context) {
	projects, err := h.biz.ListBySchool(GetUserID(c))
	if err != nil {
		if err.Error() == "permission denied" {
			Error(c, http.StatusForbidden, err.Error())
			return
		}
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	Success(c, projects)
}

func (h *ProjectHandler) ListPendingBySchool(c *gin.Context) {
	projects, err := h.biz.ListPendingBySchool(GetUserID(c))
	if err != nil {
		if err.Error() == "permission denied" {
			Error(c, http.StatusForbidden, err.Error())
			return
		}
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	Success(c, projects)
}

func (h *ProjectHandler) ApproveProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		Error(c, http.StatusBadRequest, "invalid project id")
		return
	}

	res, err := h.biz.ApproveProject(uint(id), GetUserID(c))
	if err != nil {
		switch err.Error() {
		case "permission denied":
			Error(c, http.StatusForbidden, err.Error())
		case "project not found":
			Error(c, http.StatusNotFound, err.Error())
		default:
			Error(c, http.StatusBadRequest, err.Error())
		}
		return
	}

	Success(c, res)
}

func (h *ProjectHandler) ListByCommunity(c *gin.Context) {
	projects, err := h.biz.ListByCommunity(GetUserID(c))
	if err != nil {
		if err.Error() == "permission denied" {
			Error(c, http.StatusForbidden, err.Error())
			return
		}
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	Success(c, projects)
}

func (h *ProjectHandler) ListApprovedForStudent(c *gin.Context) {
	projects, err := h.biz.ListApprovedForStudent(GetUserID(c))
	if err != nil {
		if err.Error() == "permission denied" {
			Error(c, http.StatusForbidden, err.Error())
			return
		}
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	Success(c, projects)
}
