package interfaces

import "middleware/internal/domain/models"

type TypeClpUsecase interface {
	CreateClp(typeClp *models.TypeClp) error
	SearchAllClps() (*[]models.TypeClp, error)
	SearchClpById(id uint) (*models.TypeClp, error)
}

type TypeClpRepository interface {
	CreateClp(typeClp *models.TypeClp) error
	SearchAllClps() (*[]models.TypeClp, error)
	SearchClpById(id uint) (*models.TypeClp, error)
}
