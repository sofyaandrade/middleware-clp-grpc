package usecase

import (
	"middleware/internal/domain/connection"
	"middleware/internal/domain/interfaces"
	"middleware/internal/domain/models"
	"middleware/internal/domain/security"
)

type loginUsecase struct {
	UserRepository interfaces.UserRepository
}

func NewLoginUsecase(userRepository interfaces.UserRepository) interfaces.LoginUsecase {
	return &loginUsecase{
		UserRepository: userRepository,
	}
}

func (lu *loginUsecase) SearchUserByEmail(email string) (*models.User, error) {
	return lu.UserRepository.SearchUserByEmail(email)
}

func (lu *loginUsecase) CreateAcessToken(users *models.User) (accessToken string, err error) {
	return security.CreateAcessToken(users)
}

func (lu *loginUsecase) CreateRefreshToken(users *models.User) (accessToken string, err error) {
	return security.CreateRefreshToken(users)
}

func (lu *loginUsecase) RequestToken() (repostaLogin *models.ResponseLogin, err error) {
	return connection.RequestToken()
}
