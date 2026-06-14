package interfaces

import "middleware/internal/domain/models"

type CLPUsecase interface {
	CreateClp(clp *models.CLP) error
	SearchAllClps() (*[]models.CLP, error)
	SearchClpById(id uint) (*models.CLP, error)
	SearchClpByType(typeId uint) (*[]models.CLP, error)
	UpdateClp(id uint, clp *models.CLP) error
	DeleteClp(id uint) error
	ClpsStatus() map[uint]bool
}

type CLPRepository interface {
	CreateClp(clp *models.CLP) error
	SearchAllClps() (*[]models.CLP, error)
	SearchClpById(id uint) (*models.CLP, error)
	SearchClpByType(typeId uint) (*[]models.CLP, error)
	UpdateClp(id uint, clp *models.CLP) error
	DeleteClp(id uint) error
}

type CLPReloadNotifier interface {
	RequestCLPReload(clpID uint)
}
