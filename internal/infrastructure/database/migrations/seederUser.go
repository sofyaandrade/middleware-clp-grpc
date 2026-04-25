package migrations

import (
	"fmt"
	"middleware/internal/domain/models"

	"gorm.io/gorm"
)

var listaUser = []models.User{
	{
		Name:     "adm",
		Email:    "adm",
		Password: "$2a$10$9RL6Ktc0tE3eSPiIA7x9g.AD9A6uW.CT4LcmQPU5jHL6HU5GR23UW",
	},
}

func InitializeBasicUser(db *gorm.DB) {
	var used []models.User

	err := db.Find(&used).Error

	if err != nil {
		fmt.Println("não foi possível localizar a tabela usuários: ", err)
	}
	if len(used) == 0 {
		for i := range listaUser {
			err := db.Debug().Model(&models.User{}).Create(&listaUser[i]).Error
			if err != nil {
				fmt.Println("não foi possível inserir usuários na tabela: ", err)
			}
		}
	}
}
