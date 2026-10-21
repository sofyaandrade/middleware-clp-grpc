package usecase

import (
	"middleware/internal/domain/interfaces"
	"middleware/internal/domain/models"
)

type userUsecase struct {
	usuarioRepository interfaces.UserRepository
}

func NewUserUsecase(usuarioRepository interfaces.UserRepository) interfaces.UserUsecase {
	return &userUsecase{
		usuarioRepository: usuarioRepository,
	}
}

func (uu *userUsecase) CreateUser(usuario *models.User) error {
	return uu.usuarioRepository.CreateUser(usuario)
}

func (uu *userUsecase) SearchAllUsers() (*[]models.User, error) {
	return uu.usuarioRepository.SearchAllUsers()
}

func (uu *userUsecase) SearchUserById(id uint) (*models.User, error) {
	return uu.usuarioRepository.SearchUserById(id)
}

func (uu *userUsecase) SearchUserByIdWithPassword(id uint) (*models.User, error) {
	return uu.usuarioRepository.SearchUserByIdWithPassword(id)
}

func (uu *userUsecase) SearchUserByEmail(email string) (*models.User, error) {
	return uu.usuarioRepository.SearchUserByEmail(email)
}

func (uu *userUsecase) UpdateUser(usuario *models.User, id uint) error {
	return uu.usuarioRepository.UpdateUser(usuario, id)
}

func (uu *userUsecase) DeleteUser(id uint) error {
	return uu.usuarioRepository.DeleteUser(id)
}

func (uu *userUsecase) ExistUserWithId(id uint) error {
	return uu.usuarioRepository.ExistUserWithId(id)
}
