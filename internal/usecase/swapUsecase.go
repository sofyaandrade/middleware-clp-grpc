package usecase

import (
	"middleware/internal/domain/interfaces"
	"middleware/internal/domain/models"
)

type SwapUsecase struct {
	SwapReposiotry interfaces.SwapTagRepository
}

func NewSwapUsecase(CLPRepository interfaces.SwapTagRepository) interfaces.SwapUsecase {
	return &SwapUsecase{
		SwapReposiotry: CLPRepository,
	}
}

func (su *SwapUsecase) CreateSwap(swap *models.Swap) error {
	return su.SwapReposiotry.CreateSwap(swap)
}

func (su *SwapUsecase) SearchAllSwaps() (*[]models.Swap, error) {
	return su.SwapReposiotry.SearchAllSwaps()
}

func (su *SwapUsecase) SearchSwapById(id uint) (*models.Swap, error) {
	return su.SwapReposiotry.SearchSwapById(id)
}
