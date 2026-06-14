package migrations

import (
	"fmt"
	"middleware/internal/domain/models"

	"gorm.io/gorm"
)

var listTypesOperation = []models.TypeOperation{
	{
		Model:       gorm.Model{ID: 1},
		Description: "Coil Status",
	},
	{
		Model:       gorm.Model{ID: 2},
		Description: "Input Status",
	},
	{
		Model:       gorm.Model{ID: 3},
		Description: "IHolding Registers",
	},
	{
		Model:       gorm.Model{ID: 4},
		Description: "Input Register",
	},
}

func InitializeTypesOperation(db *gorm.DB) {
	var typesOperation []models.TypeOperation

	err := db.Find(&typesOperation).Error

	if err != nil {
		fmt.Println("não foi possível localizar a tabela tipo operação: ", err)
	}
	if len(typesOperation) == 0 {
		for i := range listTypesOperation {
			err := db.Debug().Model(&models.TypeOperation{}).Create(&listTypesOperation[i]).Error
			if err != nil {
				fmt.Println("não foi possível inserir tipo operação na tabela: ", err)
			}
		}
	}
}
