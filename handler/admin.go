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
		Error(c, http.StatusNotFound, "user not found")
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
		Error(c, http.StatusNotFound, "user not found")
		return
	}

	Success(c, gin.H{"message": "user rejected"})
}
