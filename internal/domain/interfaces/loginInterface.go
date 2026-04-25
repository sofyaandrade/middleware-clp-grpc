package interfaces

import "middleware/internal/domain/models"

type LoginUsecase interface {
	SearchUserByEmail(email string) (*models.User, error)
	CreateAcessToken(user *models.User) (accessToken string, err error)
	CreateRefreshToken(user *models.User) (refreshToken string, err error)
	RequestToken() (response *models.ResponseLogin, err error)
}
