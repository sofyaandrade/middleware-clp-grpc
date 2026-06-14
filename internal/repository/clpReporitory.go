package repository

import (
	"middleware/internal/domain/interfaces"
	"middleware/internal/domain/models"

	"gorm.io/gorm"
)

type CLPRepository struct {
	Database *gorm.DB
}

func NewCLPRepository(db *gorm.DB) interfaces.CLPRepository {
	return &CLPRepository{
		Database: db,
	}
}

func (mr CLPRepository) CreateClp(clp *models.CLP) error {
	return mr.Database.Create(&clp).Error
}

func (mr CLPRepository) SearchAllClps() (*[]models.CLP, error) {
	var clps *[]models.CLP
	err := mr.Database.Model(&models.CLP{}).
		Preload("TypeClp").
		Preload("Tags.Type").
		Preload("Tags.Swap").
		Preload("Tags.OperationType").
		Order("created_at asc").
		Find(&clps).Error
	return clps, err
}

func (mr CLPRepository) SearchClpById(id uint) (*models.CLP, error) {
	var clp *models.CLP
	err := mr.Database.Model(&models.CLP{}).
		Preload("TypeClp").
		Preload("Tags.Type").
		Preload("Tags.Swap").
		Preload("Tags.OperationType").
		Order("created_at asc").
		First(&clp, id).Error
	return clp, err
}

func (mr CLPRepository) SearchClpByType(typeId uint) (*[]models.CLP, error) {
	var clps *[]models.CLP
	err := mr.Database.Model(&models.CLP{}).
		Preload("TypeClp").
		Preload("Tags.Type").
		Preload("Tags.Swap").
		Preload("Tags.OperationType").
		Where("type_clp_id = ?", typeId).
		Order("created_at asc").
		Find(&clps).Error
	return clps, err
}

func (mr CLPRepository) UpdateClp(id uint, clp *models.CLP) error {
	err := mr.Database.Model(&clp).
		Where("id = ?", id).
		Select("*").
		Updates(&clp).Error
	return err
}

func (mr CLPRepository) DeleteClp(id uint) error {
	var clp *models.CLP
	err := mr.Database.
		Preload("Tags").
		First(&clp, id).Error
	if err != nil {
		return err
	}

	err = mr.Database.Where("id_clp = ?", clp.ID).
		Delete(&models.Tag{}).Error
	if err != nil {
		return err
	}

	err = mr.Database.Delete(&clp).Error
	return err
}
