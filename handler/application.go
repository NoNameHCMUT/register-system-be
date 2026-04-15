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
