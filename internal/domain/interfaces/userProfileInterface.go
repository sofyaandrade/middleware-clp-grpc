package interfaces

import "middleware/internal/domain/models"

type UserProfileUsecase interface {
	CreateUserProfile(userProfile *models.UserProfile) error
	SearchAllUserProfiles() (*[]models.UserProfile, error)
}

type UserProfileRepository interface {
	CreateUserProfile(userProfile *models.UserProfile) error
	SearchAllUserProfiles() (*[]models.UserProfile, error)
}
