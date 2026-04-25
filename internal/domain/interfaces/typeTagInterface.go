package interfaces

import "middleware/internal/domain/models"

type TypeTagUsecase interface {
	CreateTypeTag(typeTag *models.TypeTag) error
	SearchAllTypeTags() (*[]models.TypeTag, error)
	SearchTypeTagById(id uint) (*models.TypeTag, error)
}

type TypeTagRepository interface {
	CreateTypeTag(typeTag *models.TypeTag) error
	SearchAllTypeTags() (*[]models.TypeTag, error)
	SearchTypeTagById(id uint) (*models.TypeTag, error)
}
