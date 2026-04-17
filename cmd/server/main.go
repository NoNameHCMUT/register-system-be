package main

import (
	"log"

	"register-system-be/business"
	"register-system-be/config"
	"register-system-be/email"
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
	applicationRepo := repo.NewApplicationRepo(db)

	var emailSender email.EmailSender
	if cfg.SMTPUser != "" {
		emailSender = email.NewSMTPSender(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPassword, cfg.SMTPFrom)
	} else {
		emailSender = email.NewNoopSender()
	}

	authBiz := business.NewAuthBusiness(userRepo, affiliationRepo, cfg, emailSender)
	adminBiz := business.NewAdminBusiness(userRepo, projectRepo, applicationRepo, affiliationRepo, emailSender)
	affiliationBiz := business.NewAffiliationBusiness(affiliationRepo)
	projectBiz := business.NewProjectBusiness(projectRepo, affiliationRepo, userRepo)
	applicationBiz := business.NewApplicationBusiness(applicationRepo, projectRepo, userRepo, emailSender)

	authHandler := handler.NewAuthHandler(authBiz, cfg.UploadDir)
	adminHandler := handler.NewAdminHandler(adminBiz)
	affiliationHandler := handler.NewAffiliationHandler(affiliationBiz)
	projectHandler := handler.NewProjectHandler(projectBiz, cfg.UploadDir)
	applicationHandler := handler.NewApplicationHandler(applicationBiz)

	r := router.Setup(authHandler, adminHandler, affiliationHandler, projectHandler, applicationHandler, cfg, userRepo)
	log.Printf("Server starting on :%s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal(err)
	}
}
