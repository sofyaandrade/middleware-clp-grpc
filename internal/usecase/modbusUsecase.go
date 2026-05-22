package usecase

import (
	"middleware/internal/domain/interfaces"
	"middleware/internal/domain/models"
)

type CLPUsecase struct {
	CLPReposiotry interfaces.CLPRepository
}

func NewCLPUsecase(CLPRepository interfaces.CLPRepository) interfaces.CLPUsecase {
	return &CLPUsecase{
		CLPReposiotry: CLPRepository,
	}
}

func (mr *CLPUsecase) CreateClp(clp *models.CLP) error {
	return mr.CLPReposiotry.CreateClp(clp)
}

func (mr *CLPUsecase) SearchAllClps() (*[]models.CLP, error) {
	return mr.CLPReposiotry.SearchAllClps()
}

func (mr *CLPUsecase) SearchClpById(id uint) (*models.CLP, error) {
	return mr.CLPReposiotry.SearchClpById(id)
}

func (mr *CLPUsecase) SearchClpByType(typeId uint) (*[]models.CLP, error) {
	return mr.CLPReposiotry.SearchClpByType(typeId)
}

func (mr *CLPUsecase) UpdateClp(id uint, clp *models.CLP) error {
	return mr.CLPReposiotry.UpdateClp(id, clp)
}

func (mr *CLPUsecase) DeleteClp(id uint) error {
	return mr.CLPReposiotry.DeleteClp(id)
}

func (mr *CLPUsecase) ClpsStatus() map[uint]bool {
	return nil //implementar com a conexão clp
}
