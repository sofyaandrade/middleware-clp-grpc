package migrations

import (
	"middleware/internal/domain/models"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) {
	var dbInstance *gorm.DB
	appEnv := "dev"

	if appEnv == "prod" {
		dbInstance = db
	} else {
		dbInstance = db.Debug()
	}
	dbInstance.AutoMigrate(&models.User{})
}
