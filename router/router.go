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
	applicationHandler *handler.ApplicationHandler,
	cfg *config.Config,
	userRepo repo.UserRepo,
) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS())

	r.GET(constant.HealthCheck, func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.GET(constant.Swagger, ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Static(constant.Uploads, cfg.UploadDir)

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

	usersGroup := r.Group(constant.UsersBase)
	usersGroup.Use(middleware.Auth(cfg.JWTSecret, userRepo), middleware.ActiveOnly())
	{
		usersGroup.PATCH(constant.UsersMe, authHandler.UpdateProfile)
		usersGroup.POST(constant.UsersMeAvatar, authHandler.UploadAvatar)
	}

	admin := r.Group(constant.AdminBase)
	admin.Use(middleware.Auth(cfg.JWTSecret, userRepo), middleware.ActiveOnly(), middleware.RequireRole(string(model.RoleAdmin)))
	{
		admin.GET(constant.AdminPending, adminHandler.ListPending)
		admin.POST(constant.AdminAccept, adminHandler.AcceptUser)
		admin.POST(constant.AdminReject, adminHandler.RejectUser)
		admin.GET(constant.AdminProjects, adminHandler.ListAllProjects)
		admin.GET(constant.AdminApplications, adminHandler.ListAllApplications)
		admin.POST(constant.AdminAffiliations, adminHandler.CreateAffiliation)
		admin.PATCH(constant.AdminAffiliation, adminHandler.UpdateAffiliation)
		admin.DELETE(constant.AdminAffiliation, adminHandler.DeleteAffiliation)
	}

	school := r.Group(constant.SchoolBase)
	school.Use(middleware.Auth(cfg.JWTSecret, userRepo), middleware.ActiveOnly(), middleware.RequireRole(string(model.RoleSchool)))
	{
		school.GET(constant.SchoolProjects, projectHandler.ListBySchool)
		school.GET(constant.SchoolProjectsPending, projectHandler.ListPendingBySchool)
		school.POST(constant.SchoolProjectApprove, projectHandler.ApproveProject)
		school.GET(constant.SchoolApplicants, applicationHandler.ListByProjectForSchool)
		school.POST(constant.SchoolApplicantAction, applicationHandler.SchoolAction)
	}

	community := r.Group(constant.CommunityBase)
	community.Use(middleware.Auth(cfg.JWTSecret, userRepo), middleware.ActiveOnly(), middleware.RequireRole(string(model.RoleCommunity)))
	{
		community.GET(constant.CommunityProjects, projectHandler.ListByCommunity)
		community.GET(constant.CommunityProjectApplicants, applicationHandler.ListByProjectForCommunity)
		community.POST(constant.CommunityApplicantAction, applicationHandler.CommunityAction)
	}

	student := r.Group(constant.StudentBase)
	student.Use(middleware.Auth(cfg.JWTSecret, userRepo), middleware.ActiveOnly(), middleware.RequireRole(string(model.RoleStudent)))
	{
		student.GET(constant.StudentProjects, projectHandler.ListApprovedForStudent)
		student.POST(constant.StudentApply, applicationHandler.Apply)
		student.GET(constant.StudentApplications, applicationHandler.ListByStudent)
	}

	projectsProtected := r.Group(constant.ProjectBase)
	projectsProtected.Use(middleware.Auth(cfg.JWTSecret, userRepo), middleware.ActiveOnly())
	{
		projectsProtected.POST("", projectHandler.Create)
		projectsProtected.PATCH(constant.ProjectByID, projectHandler.Update)
		projectsProtected.POST(constant.ProjectBanner, projectHandler.UploadBanner)
	}

	r.GET(constant.AffiliationBase, affiliationHandler.ListAll)

	return r
}
