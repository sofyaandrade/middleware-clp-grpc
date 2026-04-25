package interfaces

import "middleware/internal/domain/models"

type TypeOperationUsecase interface {
	CreateOperation(typeOperation *models.TypeOperation) error
	SearchAllOperations() (*[]models.TypeOperation, error)
	SearchOperationById(id uint) (*models.TypeOperation, error)
}

type TypeOperationRepository interface {
	CreateOperation(typeOperation *models.TypeOperation) error
	SearchAllOperations() (*[]models.TypeOperation, error)
	SearchOperationById(id uint) (*models.TypeOperation, error)
}
