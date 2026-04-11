package router

import (
	"register-system-be/config"
	"register-system-be/constant"
	"register-system-be/handler"
	"register-system-be/middleware"
	"register-system-be/model"

	"github.com/gin-gonic/gin"
)

func Setup(authHandler *handler.AuthHandler, cfg *config.Config) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS())

	r.GET(constant.HealthCheck, func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	auth := r.Group(constant.AuthBase)
	{
		auth.POST(constant.AuthRegister, authHandler.Register)
		auth.POST(constant.AuthLogin, authHandler.Login)
		auth.POST(constant.AuthRefresh, authHandler.Refresh)
	}

	authProtected := r.Group(constant.AuthBase)
	authProtected.Use(middleware.Auth(cfg.JWTSecret))
	{
		authProtected.GET(constant.AuthMe, authHandler.Me)
	}

	admin := r.Group("/")
	admin.Use(middleware.Auth(cfg.JWTSecret), middleware.RequireRole(string(model.RoleAdmin)))
	{
		// Add admin-only endpoints here
		_ = admin
	}

	return r
}
