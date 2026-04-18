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

// @Summary      List pending users
// @Description  Get all users awaiting approval
// @Tags         Admin
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} model.UserResponse
// @Failure      500 {object} map[string]string
// @Router       /admins/users/pending [get]
func (h *AdminHandler) ListPending(c *gin.Context) {
	users, err := h.biz.ListPending()
	if err != nil {
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	Success(c, users)
}

// @Summary      List active users
// @Description  Get all active users (is_active=true)
// @Tags         Admin
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} model.UserResponse
// @Failure      500 {object} map[string]string
// @Router       /admins/users/active [get]
func (h *AdminHandler) ListActiveUsers(c *gin.Context) {
	users, err := h.biz.ListActiveUsers()
	if err != nil {
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	Success(c, users)
}

// @Summary      Accept user
// @Description  Approve a pending user account
// @Tags         Admin
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "User ID"
// @Success      200 {object} model.UserResponse
// @Failure      404 {object} map[string]string
// @Router       /admins/users/{id}/accept [post]
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

// @Summary      Reject user
// @Description  Reject a pending user account
// @Tags         Admin
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "User ID"
// @Success      200 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /admins/users/{id}/reject [post]
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

// @Summary      List all projects
// @Description  Admin master view of all projects
// @Tags         Admin
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} model.ProjectResponse
// @Failure      500 {object} map[string]string
// @Router       /admins/projects [get]
func (h *AdminHandler) ListAllProjects(c *gin.Context) {
	projects, err := h.biz.ListAllProjects()
	if err != nil {
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	Success(c, projects)
}

// @Summary      List all applications
// @Description  Admin view of all student applications
// @Tags         Admin
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} model.StudentProjectResponse
// @Failure      500 {object} map[string]string
// @Router       /admins/applications [get]
func (h *AdminHandler) ListAllApplications(c *gin.Context) {
	apps, err := h.biz.ListAllApplications()
	if err != nil {
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	Success(c, apps)
}

// @Summary      Create affiliation
// @Description  Add a new affiliation
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body model.AffiliationCreateRequest true "Affiliation data"
// @Success      201 {object} model.AffiliationResponse
// @Failure      400 {object} map[string]string
// @Router       /admins/affiliations [post]
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

// @Summary      Update affiliation
// @Description  Update an existing affiliation
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Affiliation ID"
// @Param        body body model.AffiliationUpdateRequest true "Update fields"
// @Success      200 {object} model.AffiliationResponse
// @Failure      400 {object} map[string]string
// @Router       /admins/affiliations/{id} [patch]
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

// @Summary      Delete affiliation
// @Description  Remove an affiliation by ID
// @Tags         Admin
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Affiliation ID"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /admins/affiliations/{id} [delete]
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
