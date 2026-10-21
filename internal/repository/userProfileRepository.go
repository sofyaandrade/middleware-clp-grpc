package repository

import (
	"middleware/internal/domain/interfaces"
	"middleware/internal/domain/models"

	"gorm.io/gorm"
)

type userProfileRepository struct {
	database *gorm.DB
}

func NewUserProfileRepository(db *gorm.DB) interfaces.UserProfileRepository {
	return &userProfileRepository{
		database: db,
	}
}

func (ur *userProfileRepository) CreateUserProfile(userProfile *models.UserProfile) error {
	err := ur.database.Create(&userProfile).Error
	return err
}

func (ur *userProfileRepository) SearchAllUserProfiles() (*[]models.UserProfile, error) {
	var userProfiles []models.UserProfile
	err := ur.database.Find(&userProfiles).Error
	return &userProfiles, err
}
