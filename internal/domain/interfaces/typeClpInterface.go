package interfaces

import "middleware/internal/domain/models"

type TypeClpUsecase interface {
	CreateTypeClp(typeTypeClp *models.TypeClp) error
	SearchAllTypeClps() (*[]models.TypeClp, error)
	SearchTypeClpById(id uint) (*models.TypeClp, error)
}

type TypeClpRepository interface {
	CreateTypeClp(typeTypeClp *models.TypeClp) error
	SearchAllTypeClps() (*[]models.TypeClp, error)
	SearchTypeClpById(id uint) (*models.TypeClp, error)
}
