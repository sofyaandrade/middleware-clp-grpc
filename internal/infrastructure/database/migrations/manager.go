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
	dbInstance.AutoMigrate(&models.PermissoesAcesso{})
	dbInstance.AutoMigrate(&models.User{})
	dbInstance.AutoMigrate(&models.CLP{})
	dbInstance.AutoMigrate(&models.TypeClp{})
	dbInstance.AutoMigrate(&models.Tag{})
	dbInstance.AutoMigrate(&models.TypeTag{})
	dbInstance.AutoMigrate(&models.TypeOperation{})
	dbInstance.AutoMigrate(&models.Swap{})
	dbInstance.AutoMigrate(&models.UserProfile{})
}
