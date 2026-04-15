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

func (h *AuthHandler) Me(c *gin.Context) {
	res, err := h.biz.GetCurrentUser(GetUserID(c))
	if err != nil {
		Error(c, http.StatusNotFound, "user not found")
		return
	}

	Success(c, res)
}

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
