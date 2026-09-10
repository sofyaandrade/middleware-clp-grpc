package migrations

import (
	"fmt"
	"middleware/internal/domain/constants"
	"middleware/internal/domain/models"

	"gorm.io/gorm"
)

var listaUser = []models.User{
	{
		Name:       "adm",
		Email:      "adm",
		Permission: constants.UserProfileAdministrador,
		Password:   "$2a$10$9RL6Ktc0tE3eSPiIA7x9g.AD9A6uW.CT4LcmQPU5jHL6HU5GR23UW",
	},
}

func InitializeBasicUser(db *gorm.DB) {
	for i := range listaUser {
		var user models.User
		result := db.Where("email = ?", listaUser[i].Email).Limit(1).Find(&user)
		if result.Error != nil {
			fmt.Println("nao foi possivel localizar usuario na tabela: ", result.Error)
			continue
		}
		if result.RowsAffected == 0 {
			if err := db.Debug().Model(&models.User{}).Create(&listaUser[i]).Error; err != nil {
				fmt.Println("nao foi possivel inserir usuario na tabela: ", err)
			}
			continue
		}
		if user.Permission == "" && listaUser[i].Permission != "" {
			if err := db.Model(&user).Update("permission", listaUser[i].Permission).Error; err != nil {
				fmt.Println("nao foi possivel atualizar permissao do usuario: ", err)
			}
		}
	}
}
