package handler

import (
	"net/http"
	"strconv"

	"register-system-be/business"
	"register-system-be/model"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	biz business.AdminBusiness
}

func NewAdminHandler(biz business.AdminBusiness) *AdminHandler {
	return &AdminHandler{biz: biz}
}

func (h *AdminHandler) ListPending(c *gin.Context) {
	users, err := h.biz.ListPending()
	if err != nil {
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	Success(c, users)
}

func (h *AdminHandler) AcceptUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := h.biz.AcceptUser(uint(id))
	if err != nil {
		Error(c, http.StatusNotFound, err.Error())
		return
	}

	Success(c, user)
}

func (h *AdminHandler) RejectUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	if err := h.biz.RejectUser(uint(id)); err != nil {
		Error(c, http.StatusNotFound, err.Error())
		return
	}

	Success(c, gin.H{"message": "user rejected"})
}

func (h *AdminHandler) ListAllProjects(c *gin.Context) {
	projects, err := h.biz.ListAllProjects()
	if err != nil {
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	Success(c, projects)
}

func (h *AdminHandler) ListAllApplications(c *gin.Context) {
	apps, err := h.biz.ListAllApplications()
	if err != nil {
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	Success(c, apps)
}

func (h *AdminHandler) CreateAffiliation(c *gin.Context) {
	req, ok := Parse[model.AffiliationCreateRequest](c)
	if !ok {
		return
	}

	res, err := h.biz.CreateAffiliation(req)
	if err != nil {
		Error(c, http.StatusBadRequest, err.Error())
		return
	}

	Created(c, res)
}

func (h *AdminHandler) UpdateAffiliation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		Error(c, http.StatusBadRequest, "invalid affiliation id")
		return
	}

	req, ok := Parse[model.AffiliationUpdateRequest](c)
	if !ok {
		return
	}

	res, err := h.biz.UpdateAffiliation(uint(id), req)
	if err != nil {
		Error(c, http.StatusBadRequest, err.Error())
		return
	}

	Success(c, res)
}

func (h *AdminHandler) DeleteAffiliation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		Error(c, http.StatusBadRequest, "invalid affiliation id")
		return
	}

	if err := h.biz.DeleteAffiliation(uint(id)); err != nil {
		Error(c, http.StatusBadRequest, err.Error())
		return
	}

	Success(c, gin.H{"message": "affiliation deleted"})
}
