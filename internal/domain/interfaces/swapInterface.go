package interfaces

import "middleware/internal/domain/models"

type SwapUsecase interface {
	CreateSwap(swap *models.Swap) error
	SearchAllSwaps() (*[]models.Swap, error)
	SearchSwapById(id uint) (*models.Swap, error)
}

type SwapTagRepository interface {
	CreateSwap(swap *models.Swap) error
	SearchAllSwaps() (*[]models.Swap, error)
	SearchSwapById(id uint) (*models.Swap, error)
}
