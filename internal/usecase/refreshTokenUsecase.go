package usecase

import (
	"middleware/internal/domain/interfaces"
	"middleware/internal/domain/models"
	"middleware/internal/domain/security"
)

type refreshTokenUsecase struct {
	UserRepository interfaces.UserRepository
}

func NewRefreshTokenUsecase(userRepository interfaces.UserRepository) interfaces.RefreshTokenUsecase {
	return &refreshTokenUsecase{
		UserRepository: userRepository,
	}
}
func (rtu *refreshTokenUsecase) SearchUserById(id uint) (*models.User, error) {
	return rtu.UserRepository.SearchUserById(id)
}

func (rtu *refreshTokenUsecase) CreateAcessToken(user *models.User) (accessToken string, err error) {
	return security.CreateAcessToken(user)
}

func (rtu *refreshTokenUsecase) CreateRefreshToken(user *models.User) (refreshToken string, err error) {
	return security.CreateRefreshToken(user)
}

func (rtu *refreshTokenUsecase) ExtractIdToken(requestToken, secretKey string) (string, string, error) {
	return security.ExtractIdToken(requestToken, secretKey)
}
