package migrations

import (
	"fmt"
	"middleware/internal/domain/models"

	"gorm.io/gorm"
)

var listSwaps = []models.Swap{
	{
		Model:       gorm.Model{ID: 1},
		Description: "Sem Swap",
		OrderSwap:   "ABCD",
	},
	{
		Model:       gorm.Model{ID: 2},
		Description: "Byte swap",
		OrderSwap:   "BADC",
	},
	{
		Model:       gorm.Model{ID: 3},
		Description: "Word Swap",
		OrderSwap:   "CDAB",
	},
	{
		Model:       gorm.Model{ID: 4},
		Description: "Dword Swap", //verificar se o nome ta certo
		OrderSwap:   "DCBA",
	},
}

func InitializeSwaps(db *gorm.DB) {
	var swaps []models.Swap

	err := db.Find(&swaps).Error

	if err != nil {
		fmt.Println("não foi possível localizar a tabela swap: ", err)
	}
	if len(swaps) == 0 {
		for i := range listSwaps {
			err := db.Debug().Model(&models.Swap{}).Create(&listSwaps[i]).Error
			if err != nil {
				fmt.Println("não foi possível inserir swap na tabela: ", err)
			}
		}
	}
}
