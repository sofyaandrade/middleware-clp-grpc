package repository

import (
	"middleware/internal/domain/interfaces"
	"middleware/internal/domain/models"

	"gorm.io/gorm"
)

type SwapRepository struct {
	Database *gorm.DB
}

func NewSwapRepository(db *gorm.DB) interfaces.SwapTagRepository {
	return &SwapRepository{
		Database: db,
	}
}

func (sr SwapRepository) CreateSwap(swap *models.Swap) error {
	return sr.Database.Create(&swap).Error
}

func (sr SwapRepository) SearchAllSwaps() (*[]models.Swap, error) {
	var swaps *[]models.Swap
	err := sr.Database.Model(&models.Swap{}).Order("created_at asc").Find(&swaps).Error
	return swaps, err
}

func (sr SwapRepository) SearchSwapById(id uint) (*models.Swap, error) {
	var swap *models.Swap
	err := sr.Database.Model(&models.Swap{}).Order("created_at asc").First(&swap, id).Error
	return swap, err
}
