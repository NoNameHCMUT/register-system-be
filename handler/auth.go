package handler

import (
	"net/http"

	"register-system-be/business"
	"register-system-be/model"
	"register-system-be/upload"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	biz       business.AuthBusiness
	uploadDir string
}

func NewAuthHandler(biz business.AuthBusiness, uploadDir string) *AuthHandler {
	return &AuthHandler{biz: biz, uploadDir: uploadDir}
}

// @Summary      Register user
// @Description  Register a new user account. Account starts inactive until admin approves.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body body model.RegisterRequest true "Register data"
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
// @Description  Login with username and password. Inactive users cannot login.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body body model.LoginRequest true "Login credentials"
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
// @Description  Get new access/refresh token pair using a valid refresh token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body body model.RefreshRequest true "Refresh token"
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
// @Description  Get the profile of the currently authenticated user
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

// @Summary      Update profile
// @Description  Update current user's contact information (full_name, phone)
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body model.UserUpdateRequest true "Profile fields"
// @Success      200 {object} model.UserResponse
// @Failure      400 {object} map[string]string
// @Router       /users/me [patch]
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	req, ok := Parse[model.UserUpdateRequest](c)
	if !ok {
		return
	}

	res, err := h.biz.UpdateProfile(GetUserID(c), req)
	if err != nil {
		Error(c, http.StatusBadRequest, err.Error())
		return
	}

	Success(c, res)
}

// @Summary      Upload avatar
// @Description  Upload a profile picture for the current user
// @Tags         User
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        avatar formData file true "Avatar image"
// @Success      200 {object} model.UserResponse
// @Failure      400 {object} map[string]string
// @Router       /users/me/avatar [post]
func (h *AuthHandler) UploadAvatar(c *gin.Context) {
	file, header, err := c.Request.FormFile("avatar")
	if err != nil {
		Error(c, http.StatusBadRequest, "avatar file required")
		return
	}
	defer file.Close()

	path, err := upload.Save(file, header, "avatars", h.uploadDir)
	if err != nil {
		Error(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.biz.UpdateAvatar(GetUserID(c), path)
	if err != nil {
		Error(c, http.StatusBadRequest, err.Error())
		return
	}

	Success(c, res)
}
