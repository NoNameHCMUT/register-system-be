package handler

import (
	"net/http"
	"strconv"

	"register-system-be/business"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	biz business.AdminBusiness
}

func NewAdminHandler(biz business.AdminBusiness) *AdminHandler {
	return &AdminHandler{biz: biz}
}

// @Summary      List pending users
// @Description  Returns all users awaiting approval
// @Tags         Admin
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} model.UserResponse
// @Failure      401 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Router       /admin/users/pending [get]
func (h *AdminHandler) ListPending(c *gin.Context) {
	users, err := h.biz.ListPending()
	if err != nil {
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	Success(c, users)
}

// @Summary      Accept a pending user
// @Description  Activates a user account so they can login
// @Tags         Admin
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "User ID"
// @Success      200 {object} model.UserResponse
// @Failure      401 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /admin/users/{id}/accept [post]
func (h *AdminHandler) AcceptUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := h.biz.AcceptUser(uint(id))
	if err != nil {
		Error(c, http.StatusNotFound, "user not found")
		return
	}

	Success(c, user)
}

// @Summary      Reject a pending user
// @Description  Keeps the user account inactive
// @Tags         Admin
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "User ID"
// @Success      200 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /admin/users/{id}/reject [post]
func (h *AdminHandler) RejectUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	if err := h.biz.RejectUser(uint(id)); err != nil {
		Error(c, http.StatusNotFound, "user not found")
		return
	}

	Success(c, gin.H{"message": "user rejected"})
}
