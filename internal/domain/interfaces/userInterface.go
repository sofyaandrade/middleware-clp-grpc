package interfaces

import "middleware/internal/domain/models"

type UserUsecase interface {
	CreateUser(user *models.User) error
	SearchAllUsers() (*[]models.User, error)
	SearchUserById(id uint) (*models.User, error)
	SearchUserByIdWithPassword(id uint) (*models.User, error)
	SearchUserByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User, id uint) error
	DeleteUser(id uint) error
	ExistUserWithId(id uint) error
}

type UserRepository interface {
	CreateUser(user *models.User) error
	SearchAllUsers() (*[]models.User, error)
	SearchUserById(id uint) (*models.User, error)
	SearchUserByIdWithPassword(id uint) (*models.User, error)
	SearchUserByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User, id uint) error
	DeleteUser(id uint) error
	ExistUserWithId(id uint) error
}
