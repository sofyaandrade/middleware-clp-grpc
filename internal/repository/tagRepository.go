package repository

import (
	"middleware/internal/domain/interfaces"
	"middleware/internal/domain/models"

	"gorm.io/gorm"
)

type TagRepository struct {
	Database *gorm.DB
}

func NewTagRepository(db *gorm.DB) interfaces.TagRepository {
	return &TagRepository{
		Database: db,
	}
}

func (tr TagRepository) CreateTag(tag *models.Tag) error {
	return tr.Database.Create(&tag).Error
}

func (tr TagRepository) SearchAllTags() (*[]models.Tag, error) {
	var tags *[]models.Tag
	err := tr.Database.Model(&models.Tag{}).Preload("TypeTag").Preload("Swap").Preload("TypeOperation").Order("created_at asc").Find(&tags).Error
	return tags, err
}

func (tr TagRepository) SearchTagById(id uint) (*models.Tag, error) {
	var tag *models.Tag
	err := tr.Database.Model(&models.Tag{}).Preload("TypeTag").Preload("Swap").Preload("TypeOperation").Order("created_at asc").First(&tag, id).Error
	return tag, err
}

func (tr TagRepository) UpdateTag(id uint, tag *models.Tag) error {
	err := tr.Database.Model(&tag).Where("id = ?", id).Select("*").Updates(&tag).Error
	return err
}

func (tr TagRepository) DeleteTag(id uint) error {
	var tag *models.Tag
	err := tr.Database.Delete(&tag).Error
	return err
}

func (tr TagRepository) ExistTagWithId(id uint) error {
	err := tr.Database.Model(&models.Tag{}).Preload("TypeTag").Preload("Swap").Preload("TypeOperation").Order("created_at asc").First(id).Error
	return err
}
