package repository

import (
	"middleware/internal/domain/interfaces"
	"middleware/internal/domain/models"

	"gorm.io/gorm"
)

type TypeClpRepository struct {
	Database *gorm.DB
}

func NewTypeClpRepository(db *gorm.DB) interfaces.TypeClpRepository {
	return &TypeClpRepository{
		Database: db,
	}
}

func (ttr TypeClpRepository) CreateTypeClp(typeClp *models.TypeClp) error {
	return ttr.Database.Create(&typeClp).Error
}

func (ttr TypeClpRepository) SearchAllTypeClps() (*[]models.TypeClp, error) {
	var typeClps *[]models.TypeClp
	err := ttr.Database.
		Model(&models.TypeClp{}).
		Order("created_at asc").
		Find(&typeClps).Error
	return typeClps, err
}

func (ttr TypeClpRepository) SearchTypeClpById(id uint) (*models.TypeClp, error) {
	var typeClp *models.TypeClp
	err := ttr.Database.
		Model(&models.TypeClp{}).
		Order("created_at asc").
		First(&typeClp, id).Error
	return typeClp, err
}
