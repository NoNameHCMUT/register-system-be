package repo

import (
	"fmt"
	"log"

	"register-system-be/config"
	"register-system-be/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}
	return db
}

func Migrate(db *gorm.DB) {
	if err := db.AutoMigrate(&model.Affiliation{}, &model.User{}, &model.Project{}, &model.StudentProject{}); err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}
}
