package controllers

import (
	"middleware/internal/domain/conversion"
	"middleware/internal/domain/interfaces"
	"middleware/internal/domain/models"
	"middleware/internal/domain/security"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserUsecase interfaces.UserUsecase
}

type changePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
	SenhaAtual      string `json:"senha_atual"`
	SenhaNova       string `json:"senha_nova"`
	SenhaAtualCamel string `json:"senhaAtual"`
	SenhaNovaCamel  string `json:"senhaNova"`
}

func (request changePasswordRequest) currentPassword() string {
	if request.CurrentPassword != "" {
		return request.CurrentPassword
	}
	if request.SenhaAtual != "" {
		return request.SenhaAtual
	}
	return request.SenhaAtualCamel
}

func (request changePasswordRequest) newPassword() string {
	if request.NewPassword != "" {
		return request.NewPassword
	}
	if request.SenhaNova != "" {
		return request.SenhaNova
	}
	return request.SenhaNovaCamel
}

var erroUserNaoEncontrado = "record not found"

func (uc *UserController) NewUser(c *gin.Context) {

	var users models.User

	if err := c.ShouldBindJSON(&users); err != nil {
		c.JSON(http.StatusBadRequest, "Dados da requisição inválidos.")
		return
	}

	user, err := uc.UserUsecase.SearchUserByEmail(users.Email)

	if err != nil && err.Error() != erroUserNaoEncontrado {
		c.JSON(http.StatusBadRequest, "Não foi possível verificar se o e-mail já está cadastrado.")
		return
	}

	if err == nil && user.Email != "" {
		c.JSON(http.StatusBadRequest, "Este e-mail já está cadastrado.")
		return
	}

	hash, err := security.HashSenha(users.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Não foi possível proteger a senha.")
		return
	}

	users.Password = string(hash)

	err = uc.UserUsecase.CreateUser(&users)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Não foi possível criar o usuário.")
		return
	}

	c.JSON(http.StatusOK, "Cadastro realizado com sucesso.")
}

func (uc *UserController) SearchAllUsers(c *gin.Context) {

	users, err := uc.UserUsecase.SearchAllUsers()
	if err != nil {
		c.JSON(http.StatusBadRequest, "Não foi possível buscar os usuários.")
		return
	}
	if users == nil {
		c.JSON(http.StatusBadRequest, "Nenhum usuário encontrado.")
		return
	}

	c.JSON(200, users)
}

func (uc *UserController) SearchUserById(c *gin.Context) {

	id := c.Params.ByName("id")

	user, err := uc.UserUsecase.SearchUserById(conversion.StringToUint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, "Não foi possível buscar o usuário.")
		return
	}

	if user.ID == 0 {
		c.JSON(http.StatusNotFound, "Usuário não encontrado.")
		return
	}

	c.JSON(http.StatusOK, user)
}

func (uc *UserController) UpdateUser(c *gin.Context) {
	var user *models.User

	id := c.Params.ByName("id")

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, "Dados da requisição inválidos.")
		return
	}

	err := uc.UserUsecase.ExistUserWithId(conversion.StringToUint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, "Usuário não encontrado.")
		return
	}

	if user.Password != "" {
		hash, err := security.HashSenha(user.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, "Não foi possível proteger a senha.")
			return
		}

		user.Password = string(hash)
	}

	err = uc.UserUsecase.UpdateUser(user, conversion.StringToUint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, "Não foi possível atualizar o usuário.")
		return
	}

	c.JSON(http.StatusOK, "Atualização realizada com sucesso.")

}

func (uc *UserController) ChangePassword(c *gin.Context) {
	var request changePasswordRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, "Dados da requisição inválidos.")
		return
	}

	currentPassword := request.currentPassword()
	newPassword := request.newPassword()
	if currentPassword == "" || newPassword == "" {
		c.JSON(http.StatusBadRequest, "Informe a senha atual e a nova senha.")
		return
	}

	id := c.GetString("id")
	if id == "" {
		c.JSON(http.StatusUnauthorized, "Usuário não autenticado.")
		return
	}

	user, err := uc.UserUsecase.SearchUserByIdWithPassword(conversion.StringToUint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, "Usuário não encontrado.")
		return
	}

	err = security.VerificaSenha(user.Password, currentPassword)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "A senha atual está incorreta.")
		return
	}

	hash, err := security.HashSenha(newPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Não foi possível proteger a senha.")
		return
	}

	err = uc.UserUsecase.UpdateUser(&models.User{Password: string(hash)}, user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Não foi possível atualizar a senha.")
		return
	}

	c.JSON(http.StatusOK, "Senha atualizada com sucesso.")
}

func (uc *UserController) DeleteUser(c *gin.Context) {
	id := c.Params.ByName("id")

	err := uc.UserUsecase.ExistUserWithId(conversion.StringToUint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, "Usuário não encontrado.")
		return
	}

	err = uc.UserUsecase.DeleteUser(conversion.StringToUint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, "Não foi possível excluir o usuário.")
		return
	}

	c.JSON(http.StatusOK, "Exclusão realizada com sucesso.")
}
