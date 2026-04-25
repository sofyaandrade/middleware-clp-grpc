package usecase

import (
	"middleware/internal/domain/interfaces"
	"middleware/internal/domain/models"
)

type TypeTagUsecase struct {
	TypeTagReposiotry interfaces.TypeTagRepository
}

func NewTypeTagUsecase(CLPRepository interfaces.TypeTagRepository) interfaces.TypeTagUsecase {
	return &TypeTagUsecase{
		TypeTagReposiotry: CLPRepository,
	}
}

func (ttu *TypeTagUsecase) CreateTypeTag(typeTag *models.TypeTag) error {
	return ttu.TypeTagReposiotry.CreateTypeTag(typeTag)
}

func (ttu *TypeTagUsecase) SearchAllTypeTags() (*[]models.TypeTag, error) {
	return ttu.TypeTagReposiotry.SearchAllTypeTags()
}

func (ttu *TypeTagUsecase) SearchTypeTagById(id uint) (*models.TypeTag, error) {
	return ttu.TypeTagReposiotry.SearchTypeTagById(id)
}
