package router

import (
	"register-system-be/config"
	"register-system-be/constant"
	"register-system-be/handler"
	"register-system-be/middleware"
	"register-system-be/model"
	"register-system-be/repo"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Setup(
	authHandler *handler.AuthHandler,
	adminHandler *handler.AdminHandler,
	affiliationHandler *handler.AffiliationHandler,
	projectHandler *handler.ProjectHandler,
	userRepo repo.UserRepo,
	cfg *config.Config,
) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS())

	r.GET(constant.HealthCheck, func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.GET(constant.Swagger, ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := r.Group(constant.AuthBase)
	{
		auth.POST(constant.AuthRegister, authHandler.Register)
		auth.POST(constant.AuthLogin, authHandler.Login)
		auth.POST(constant.AuthRefresh, authHandler.Refresh)
	}

	authProtected := r.Group(constant.AuthBase)
	authProtected.Use(middleware.Auth(cfg.JWTSecret, userRepo))
	{
		authProtected.GET(constant.AuthMe, authHandler.Me)
	}

	admin := r.Group(constant.AdminBase)
	admin.Use(middleware.Auth(cfg.JWTSecret, userRepo), middleware.RequireRole(string(model.RoleAdmin)))
	{
		admin.GET(constant.AdminPending, adminHandler.ListPending)
		admin.POST(constant.AdminAccept, adminHandler.AcceptUser)
		admin.POST(constant.AdminReject, adminHandler.RejectUser)
	}

	r.GET(constant.AffiliationBase, affiliationHandler.ListAll)

	projects := r.Group(constant.ProjectBase)
	projects.Use(middleware.Auth(cfg.JWTSecret, userRepo), middleware.RequireRole(string(model.RoleCommunity)))
	{
		projects.POST(constant.ProjectCreate, projectHandler.Create)
	}

	return r
}
