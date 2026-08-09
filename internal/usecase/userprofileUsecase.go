package usecase

import (
	"middleware/internal/domain/interfaces"
	"middleware/internal/domain/models"
)

type userProfileUsecase struct {
	usuarioRepository interfaces.UserProfileRepository
}

func NewUserProfileUsecase(usuarioRepository interfaces.UserProfileRepository) interfaces.UserProfileUsecase {
	return &userProfileUsecase{
		usuarioRepository: usuarioRepository,
	}
}

func (uu *userProfileUsecase) CreateUserProfile(usuario *models.UserProfile) error {
	return uu.usuarioRepository.CreateUserProfile(usuario)
}

func (uu *userProfileUsecase) SearchAllUserProfiles() (*[]models.UserProfile, error) {
	return uu.usuarioRepository.SearchAllUserProfiles()
}
