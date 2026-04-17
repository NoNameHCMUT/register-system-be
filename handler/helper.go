package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, gin.H{"data": data})
}

func Error(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"error": msg})
}

func Parse[T any](c *gin.Context) (*T, bool) {
	var req T
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, err.Error())
		return nil, false
	}
	return &req, true
}

func GetUserID(c *gin.Context) uint {
	return c.MustGet("user_id").(uint)
}

func GetRole(c *gin.Context) string {
	return c.MustGet("role").(string)
}

func GetIsActive(c *gin.Context) bool {
	return c.GetBool("is_active")
}
