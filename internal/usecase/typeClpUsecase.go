package usecase

import (
	"middleware/internal/domain/interfaces"
	"middleware/internal/domain/models"
)

type TypeClpUsecase struct {
	TypeClpReposiotry interfaces.TypeClpRepository
}

func NewTypeClpUsecase(CLPRepository interfaces.TypeClpRepository) interfaces.TypeClpUsecase {
	return &TypeClpUsecase{
		TypeClpReposiotry: CLPRepository,
	}
}

func (ttu *TypeClpUsecase) CreateTypeClp(typeClp *models.TypeClp) error {
	return ttu.TypeClpReposiotry.CreateTypeClp(typeClp)
}

func (ttu *TypeClpUsecase) SearchAllTypeClps() (*[]models.TypeClp, error) {
	return ttu.TypeClpReposiotry.SearchAllTypeClps()
}

func (ttu *TypeClpUsecase) SearchTypeClpById(id uint) (*models.TypeClp, error) {
	return ttu.TypeClpReposiotry.SearchTypeClpById(id)
}
