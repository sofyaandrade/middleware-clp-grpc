package migrations

import (
	"fmt"
	"middleware/internal/domain/constants"
	"middleware/internal/domain/models"

	"gorm.io/gorm"
)

var listaUserProfile = []models.UserProfile{
	{
		Description: constants.UserProfileAdministrador,
	},
	{
		Description: constants.UserProfileOperador,
	},
	{
		Description: constants.UserProfileConsumidor,
	},
}

func InitializeBasicUserProfile(db *gorm.DB) {
	for i := range listaUserProfile {
		var userProfile models.UserProfile
		result := db.Where("description = ?", listaUserProfile[i].Description).Limit(1).Find(&userProfile)
		if result.Error != nil {
			fmt.Println("nao foi possivel localizar perfil de usuario na tabela: ", result.Error)
			continue
		}
		if result.RowsAffected == 0 {
			if err := db.Debug().Model(&models.UserProfile{}).Create(&listaUserProfile[i]).Error; err != nil {
				fmt.Println("nao foi possivel inserir perfil de usuario na tabela: ", err)
			}
		}
	}
}
