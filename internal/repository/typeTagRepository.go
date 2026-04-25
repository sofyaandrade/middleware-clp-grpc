package repository

import (
	"middleware/internal/domain/interfaces"
	"middleware/internal/domain/models"

	"gorm.io/gorm"
)

type TypeTagRepository struct {
	Database *gorm.DB
}

func NewTypeTagRepository(db *gorm.DB) interfaces.TypeTagRepository {
	return &TypeTagRepository{
		Database: db,
	}
}

func (ttr TypeTagRepository) CreateTypeTag(typeTag *models.TypeTag) error {
	return ttr.Database.Create(&typeTag).Error
}

func (ttr TypeTagRepository) SearchAllTypeTags() (*[]models.TypeTag, error) {
	var typeTags *[]models.TypeTag
	err := ttr.Database.Model(&models.TypeTag{}).Order("created_at asc").Find(&typeTags).Error
	return typeTags, err
}

func (ttr TypeTagRepository) SearchTypeTagById(id uint) (*models.TypeTag, error) {
	var typeTag *models.TypeTag
	err := ttr.Database.Model(&models.TypeTag{}).Order("created_at asc").First(&typeTag, id).Error
	return typeTag, err
}
