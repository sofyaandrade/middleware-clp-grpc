package repository

import (
	"middleware/internal/domain/interfaces"
	"middleware/internal/domain/models"

	"gorm.io/gorm"
)

type OperationTypeRepository struct {
	Database *gorm.DB
}

func NewOperationTypeRepository(db *gorm.DB) interfaces.TypeOperationRepository {
	return &OperationTypeRepository{
		Database: db,
	}
}

func (otr OperationTypeRepository) CreateOperation(operationType *models.TypeOperation) error {
	return otr.Database.Create(&operationType).Error
}

func (otr OperationTypeRepository) SearchAllOperations() (*[]models.TypeOperation, error) {
	var operationTypes *[]models.TypeOperation
	err := otr.Database.Model(&models.TypeOperation{}).Order("created_at asc").Find(&operationTypes).Error
	return operationTypes, err
}

func (otr OperationTypeRepository) SearchOperationById(id uint) (*models.TypeOperation, error) {
	var operationType *models.TypeOperation
	err := otr.Database.Model(&models.TypeOperation{}).Order("created_at asc").First(&operationType, id).Error
	return operationType, err
}
