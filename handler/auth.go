package handler

import (
	"net/http"

	"register-system-be/business"
	"register-system-be/model"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	biz business.AuthBusiness
}

func NewAuthHandler(biz business.AuthBusiness) *AuthHandler {
	return &AuthHandler{biz: biz}
}

// @Summary      Register a new user
// @Description  Creates an inactive account pending admin approval
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body body model.RegisterRequest true "Register request"
// @Success      201 {object} model.RegisterResponse
// @Failure      400 {object} map[string]string
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	req, ok := Parse[model.RegisterRequest](c)
	if !ok {
		return
	}

	res, err := h.biz.Register(req)
	if err != nil {
		Error(c, http.StatusBadRequest, err.Error())
		return
	}

	Created(c, res)
}

// @Summary      Login
// @Description  Authenticates a user and returns JWT tokens. Only active accounts can login.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body body model.LoginRequest true "Login request"
// @Success      200 {object} model.AuthResponse
// @Failure      401 {object} map[string]string
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	req, ok := Parse[model.LoginRequest](c)
	if !ok {
		return
	}

	res, err := h.biz.Login(req)
	if err != nil {
		Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	Success(c, res)
}

// @Summary      Refresh tokens
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body body model.RefreshRequest true "Refresh request"
// @Success      200 {object} model.AuthResponse
// @Failure      401 {object} map[string]string
// @Router       /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	req, ok := Parse[model.RefreshRequest](c)
	if !ok {
		return
	}

	res, err := h.biz.Refresh(req)
	if err != nil {
		Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	Success(c, res)
}

// @Summary      Get current user
// @Tags         Auth
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} model.UserResponse
// @Failure      401 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	res, err := h.biz.GetCurrentUser(GetUserID(c))
	if err != nil {
		Error(c, http.StatusNotFound, "user not found")
		return
	}

	Success(c, res)
}
