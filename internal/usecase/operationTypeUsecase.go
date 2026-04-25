package usecase

import (
	"middleware/internal/domain/interfaces"
	"middleware/internal/domain/models"
)

type OperationTypeUsecase struct {
	OperationTypeReposiotry interfaces.TypeOperationRepository
}

func NewOperationTypeUsecase(CLPRepository interfaces.TypeOperationRepository) interfaces.TypeOperationUsecase {
	return &OperationTypeUsecase{
		OperationTypeReposiotry: CLPRepository,
	}
}

func (otu *OperationTypeUsecase) CreateOperation(clp *models.TypeOperation) error {
	return otu.OperationTypeReposiotry.CreateOperation(clp)
}

func (otu *OperationTypeUsecase) SearchAllOperations() (*[]models.TypeOperation, error) {
	return otu.OperationTypeReposiotry.SearchAllOperations()
}

func (otu *OperationTypeUsecase) SearchOperationById(id uint) (*models.TypeOperation, error) {
	return otu.OperationTypeReposiotry.SearchOperationById(id)
}
