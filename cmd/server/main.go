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

	authBiz := business.NewAuthBusiness(userRepo, affiliationRepo, cfg)
	adminBiz := business.NewAdminBusiness(userRepo)

	authHandler := handler.NewAuthHandler(authBiz)
	adminHandler := handler.NewAdminHandler(adminBiz)

	r := router.Setup(authHandler, adminHandler, cfg)
	log.Printf("Server starting on :%s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal(err)
	}
}
