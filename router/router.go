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
	cfg *config.Config,
	userRepo repo.UserRepo,
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
	admin.Use(middleware.Auth(cfg.JWTSecret, userRepo), middleware.ActiveOnly(), middleware.RequireRole(string(model.RoleAdmin)))
	{
		admin.GET(constant.AdminPending, adminHandler.ListPending)
		admin.POST(constant.AdminAccept, adminHandler.AcceptUser)
		admin.POST(constant.AdminReject, adminHandler.RejectUser)
	}

	r.GET(constant.AffiliationBase, affiliationHandler.ListAll)

	project := r.Group(constant.ProjectBase)
	project.Use(middleware.Auth(cfg.JWTSecret, userRepo), middleware.ActiveOnly())
	{
		project.GET(constant.ProjectMyList, projectHandler.GetMyProjects)
		project.POST("", projectHandler.Create)
		project.PATCH(constant.ProjectByID, projectHandler.Update)
	}
	return r
}
