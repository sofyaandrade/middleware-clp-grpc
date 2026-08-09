package migrations

import (
	"fmt"
	"middleware/internal/domain/models"

	"gorm.io/gorm"
)

var listaUserProfile = []models.UserProfile{
	{
		Description: "ADMINISTRADOR",
	},
	{
		Description: "OPERADOR",
	},
}

func InitializeBasicUserProfile(db *gorm.DB) {
	var used []models.UserProfile

	err := db.Find(&used).Error

	if err != nil {
		fmt.Println("não foi possível localizar a tabela perfis de usuário: ", err)
	}
	if len(used) == 0 {
		for i := range listaUserProfile {
			err := db.Debug().Model(&models.UserProfile{}).Create(&listaUserProfile[i]).Error
			if err != nil {
				fmt.Println("não foi possível inserir perfis de usuário na tabela: ", err)
			}
		}
	}
}
