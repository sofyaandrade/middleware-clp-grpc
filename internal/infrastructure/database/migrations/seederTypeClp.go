package migrations

import (
	"fmt"
	"middleware/internal/domain/models"

	"gorm.io/gorm"
)

var listTypesClp = []models.TypeClp{
	{
		Model:       gorm.Model{ID: 1},
		Description: "Modbus Master",
	},
	{
		Model:       gorm.Model{ID: 2},
		Description: "Modbus Slave",
	},
}

func InitializeTypesClp(db *gorm.DB) {
	var typesClp []models.TypeClp

	err := db.Find(&typesClp).Error

	if err != nil {
		fmt.Println("não foi possível localizar a tabela tipo clp: ", err)
	}
	if len(typesClp) == 0 {
		for i := range listTypesClp {
			err := db.Debug().Model(&models.TypeClp{}).Create(&listTypesClp[i]).Error
			if err != nil {
				fmt.Println("não foi possível inserir tipo clp na tabela: ", err)
			}
		}
	}
}
