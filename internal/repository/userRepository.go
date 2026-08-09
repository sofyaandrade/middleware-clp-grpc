package repository

import (
	"middleware/internal/domain/interfaces"
	"middleware/internal/domain/models"

	"gorm.io/gorm"
)

type userRepository struct {
	database *gorm.DB
}

func NewUserRepository(db *gorm.DB) interfaces.UserRepository {
	return &userRepository{
		database: db,
	}
}

func (ur *userRepository) CreateUser(user *models.User) error {
	err := ur.database.Create(&user).Error
	return err
}

func (ur *userRepository) SearchAllUsers() (*[]models.User, error) {
	var user *[]models.User
	err := ur.database.Select("id", "name", "email", "phone", "permission").Find(&user).Error
	return user, err
}

func (ur *userRepository) SearchUserById(id uint) (*models.User, error) {
	var user *models.User
	err := ur.database.Select("id", "name", "email", "phone", "permission").First(&user, id).Error
	return user, err
}

func (ur *userRepository) SearchUserByEmail(email string) (*models.User, error) {
	var user *models.User
	err := ur.database.Where(models.User{Email: email}).First(&user).Error
	return user, err
}

func (ur *userRepository) UpdateUser(user *models.User, id uint) error {
	err := ur.database.Model(&user).Where("id = ?", id).Updates(user).Error
	return err
}

func (ur *userRepository) DeleteUser(id uint) error {
	var user *models.User
	err := ur.database.Delete(&user, id).Error
	return err
}

func (ur *userRepository) ExistUserWithId(id uint) error {
	var user *models.User
	err := ur.database.First(&user, id).Error
	return err
}
