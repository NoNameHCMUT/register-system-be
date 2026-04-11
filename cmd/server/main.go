package main

import (
	"log"

	"register-system-be/business"
	"register-system-be/config"
	"register-system-be/handler"
	"register-system-be/repo"
	"register-system-be/router"
)

func main() {
	cfg := config.Load()
	db := repo.Connect(cfg)
	repo.Migrate(db)

	userRepo := repo.NewUserRepo(db)
	tokenRepo := repo.NewTokenRepo(db)
	authBiz := business.NewAuthBusiness(userRepo, tokenRepo, cfg)
	authHandler := handler.NewAuthHandler(authBiz)

	r := router.Setup(authHandler, cfg)
	log.Printf("Server starting on :%s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal(err)
	}
}
