package interfaces

import "middleware/internal/domain/models"

type RefreshTokenUsecase interface {
	SearchUserById(id uint) (*models.User, error)
	CreateAcessToken(user *models.User) (accessToken string, err error)
	CreateRefreshToken(user *models.User) (refreshToken string, err error)
	ExtractIdToken(requestToken, secretKey string) (userId string, role string, err error)
}
