package migrations

import (
	"fmt"
	"middleware/internal/domain/models"

	"gorm.io/gorm"
)

var listTypesTag = []models.TypeTag{
	{
		Model:       gorm.Model{ID: 1},
		Description: "Real",
	},
	{
		Model:       gorm.Model{ID: 2},
		Description: "Int",
	},
	{
		Model:       gorm.Model{ID: 3},
		Description: "Bool",
	},
	{
		Model:       gorm.Model{ID: 4},
		Description: "Dword",
	},
}

func InitializeTypesTag(db *gorm.DB) {
	var typesTag []models.TypeTag

	err := db.Find(&typesTag).Error

	if err != nil {
		fmt.Println("não foi possível localizar a tabela tipo tag: ", err)
	}
	if len(typesTag) == 0 {
		for i := range listTypesTag {
			err := db.Debug().Model(&models.TypeTag{}).Create(&listTypesTag[i]).Error
			if err != nil {
				fmt.Println("não foi possível inserir tipo tag na tabela: ", err)
			}
		}
	}
}
