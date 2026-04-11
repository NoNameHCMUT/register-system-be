package main

import (
	"log"

	"register-system-be/business"
	"register-system-be/config"
	"register-system-be/handler"
	"register-system-be/repo"
	"register-system-be/router"

	_ "register-system-be/docs"
)

// @title           Register System API
// @version         1.0
// @description     Backend API for the register system with JWT auth and role-based access.
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg := config.Load()
	db := repo.Connect(cfg)
	repo.Migrate(db)

	userRepo := repo.NewUserRepo(db)
	affiliationRepo := repo.NewAffiliationRepo(db)
	projectRepo := repo.NewProjectRepo(db)

	authBiz := business.NewAuthBusiness(userRepo, affiliationRepo, cfg)
	adminBiz := business.NewAdminBusiness(userRepo)
	affiliationBiz := business.NewAffiliationBusiness(affiliationRepo)
	projectBiz := business.NewProjectBusiness(projectRepo)

	authHandler := handler.NewAuthHandler(authBiz)
	adminHandler := handler.NewAdminHandler(adminBiz)
	affiliationHandler := handler.NewAffiliationHandler(affiliationBiz)
	projectHandler := handler.NewProjectHandler(projectBiz, userRepo)

	r := router.Setup(authHandler, adminHandler, affiliationHandler, projectHandler, cfg, userRepo)
	log.Printf("Server starting on :%s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal(err)
	}
}