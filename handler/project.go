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

// @Summary      Create project
// @Description  Create a new project (school or community only)
// @Tags         Project
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body model.ProjectCreateRequest true "Project data"
// @Success      201 {object} model.ProjectResponse
// @Failure      400 {object} map[string]string
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

// @Summary      Update project
// @Description  Update a project owned by the current user
// @Tags         Project
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Project ID"
// @Param        body body model.ProjectUpdateRequest true "Update fields"
// @Success      200 {object} model.ProjectResponse
// @Failure      400 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /projects/{id} [patch]
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

// @Summary      Upload project banner
// @Description  Upload a banner image for a project
// @Tags         Project
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Project ID"
// @Param        banner formData file true "Banner image"
// @Success      200 {object} model.ProjectResponse
// @Failure      400 {object} map[string]string
// @Router       /projects/{id}/banner [post]
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

// @Summary      List school projects
// @Description  Get all projects belonging to the school's affiliation
// @Tags         School
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} model.ProjectResponse
// @Failure      403 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /schools/projects [get]
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

// @Summary      List pending school projects
// @Description  Get projects pending approval in the school's affiliation
// @Tags         School
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} model.ProjectResponse
// @Failure      403 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /schools/projects/pending [get]
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

// @Summary      Approve project
// @Description  School approves a pending project
// @Tags         School
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Project ID"
// @Success      200 {object} model.ProjectResponse
// @Failure      400 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /schools/projects/{id}/approve [post]
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

// @Summary      List community projects
// @Description  Get all approved projects for the community's affiliation
// @Tags         Community
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} model.ProjectResponse
// @Failure      403 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /community/projects [get]
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

// @Summary      List approved projects for student
// @Description  Get approved projects available in the student's affiliation
// @Tags         Student
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} model.ProjectResponse
// @Failure      403 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /students/projects [get]
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
