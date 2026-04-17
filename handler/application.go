package handler

import (
	"net/http"
	"strconv"

	"register-system-be/business"
	"register-system-be/model"

	"github.com/gin-gonic/gin"
)

type ApplicationHandler struct {
	biz business.ApplicationBusiness
}

func NewApplicationHandler(biz business.ApplicationBusiness) *ApplicationHandler {
	return &ApplicationHandler{biz: biz}
}

// @Summary      Apply to project
// @Description  Student applies to a project by its ID
// @Tags         Student
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Project ID"
// @Success      201 {object} model.StudentProjectResponse
// @Failure      400 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /students/projects/{id}/apply [post]
func (h *ApplicationHandler) Apply(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		Error(c, http.StatusBadRequest, "invalid project id")
		return
	}

	res, err := h.biz.Apply(GetUserID(c), uint(projectID))
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

	Created(c, res)
}

// @Summary      List student applications
// @Description  Get all applications submitted by the current student
// @Tags         Student
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} model.StudentProjectResponse
// @Failure      403 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /students/applications [get]
func (h *ApplicationHandler) ListByStudent(c *gin.Context) {
	apps, err := h.biz.ListByStudent(GetUserID(c))
	if err != nil {
		if err.Error() == "permission denied" {
			Error(c, http.StatusForbidden, err.Error())
			return
		}
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	Success(c, apps)
}

// @Summary      List applicants for school
// @Description  Get all applications for a project (school view)
// @Tags         School
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Project ID"
// @Success      200 {array} model.StudentProjectResponse
// @Failure      400 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Router       /schools/projects/{id}/applicants [get]
func (h *ApplicationHandler) ListByProjectForSchool(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		Error(c, http.StatusBadRequest, "invalid project id")
		return
	}

	apps, err := h.biz.ListByProjectForSchool(GetUserID(c), uint(projectID))
	if err != nil {
		if err.Error() == "permission denied" {
			Error(c, http.StatusForbidden, err.Error())
			return
		}
		Error(c, http.StatusBadRequest, err.Error())
		return
	}
	Success(c, apps)
}

// @Summary      School action on applications
// @Description  Batch accept or reject applications as a school
// @Tags         School
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body model.ApplicationActionRequest true "Action data"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Router       /schools/applicants/action [post]
func (h *ApplicationHandler) SchoolAction(c *gin.Context) {
	req, ok := Parse[model.ApplicationActionRequest](c)
	if !ok {
		return
	}

	if err := h.biz.SchoolAction(GetUserID(c), req); err != nil {
		if err.Error() == "permission denied" {
			Error(c, http.StatusForbidden, err.Error())
			return
		}
		Error(c, http.StatusBadRequest, err.Error())
		return
	}

	Success(c, gin.H{"message": "action processed"})
}

// @Summary      List applicants for community
// @Description  Get all applications for a project (community view)
// @Tags         Community
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Project ID"
// @Success      200 {array} model.StudentProjectResponse
// @Failure      400 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Router       /community/projects/{id}/applicants [get]
func (h *ApplicationHandler) ListByProjectForCommunity(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		Error(c, http.StatusBadRequest, "invalid project id")
		return
	}

	apps, err := h.biz.ListByProjectForCommunity(GetUserID(c), uint(projectID))
	if err != nil {
		if err.Error() == "permission denied" {
			Error(c, http.StatusForbidden, err.Error())
			return
		}
		Error(c, http.StatusBadRequest, err.Error())
		return
	}
	Success(c, apps)
}

// @Summary      Community action on applications
// @Description  Batch accept or reject applications as a community
// @Tags         Community
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body model.ApplicationActionRequest true "Action data"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Router       /community/applicants/action [post]
func (h *ApplicationHandler) CommunityAction(c *gin.Context) {
	req, ok := Parse[model.ApplicationActionRequest](c)
	if !ok {
		return
	}

	if err := h.biz.CommunityAction(GetUserID(c), req); err != nil {
		if err.Error() == "permission denied" {
			Error(c, http.StatusForbidden, err.Error())
			return
		}
		Error(c, http.StatusBadRequest, err.Error())
		return
	}

	Success(c, gin.H{"message": "action processed"})
}
