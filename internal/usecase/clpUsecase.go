package usecase

import (
	"middleware/internal/domain/interfaces"
	"middleware/internal/domain/models"
	"middleware/internal/infrastructure/jobs"
)

type CLPUsecase struct {
	CLPReposiotry  interfaces.CLPRepository
	ReloadNotifier interfaces.CLPReloadNotifier
}

func NewCLPUsecase(CLPRepository interfaces.CLPRepository, reloadNotifier interfaces.CLPReloadNotifier) interfaces.CLPUsecase {
	return &CLPUsecase{
		CLPReposiotry:  CLPRepository,
		ReloadNotifier: reloadNotifier,
	}
}

func (mr *CLPUsecase) CreateClp(clp *models.CLP) error {
	if err := mr.CLPReposiotry.CreateClp(clp); err != nil {
		return err
	}

	mr.requestCLPReload(clp.ID)
	return nil
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
	if err := mr.CLPReposiotry.UpdateClp(id, clp); err != nil {
		return err
	}

	mr.requestCLPReload(id)
	return nil
}

func (mr *CLPUsecase) DeleteClp(id uint) error {
	if err := mr.CLPReposiotry.DeleteClp(id); err != nil {
		return err
	}

	mr.requestCLPReload(id)
	return nil
}

func (mr *CLPUsecase) ClpsStatus() map[uint]bool {
	return jobs.ReadAllClpsStatus()
}

func (mr *CLPUsecase) requestCLPReload(clpID uint) {
	if mr.ReloadNotifier != nil {
		mr.ReloadNotifier.RequestCLPReload(clpID)
	}
}
