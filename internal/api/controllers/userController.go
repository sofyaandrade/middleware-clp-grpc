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

var erroUserNaoEncontrado = "record not found"

func (uc *UserController) NewUser(c *gin.Context) {

	var users models.User

	if err := c.ShouldBindJSON(&users); err != nil {
		c.JSON(http.StatusBadRequest, "error.json")
		return
	}

	user, err := uc.UserUsecase.SearchUserByEmail(users.Email)

	if err != nil && err.Error() != erroUserNaoEncontrado {
		c.JSON(http.StatusBadRequest, "not.found")
		return
	}

	if err == nil && user.Email != "" {
		c.JSON(http.StatusBadRequest, "email.alredy.used")
		return
	}

	hash, err := security.HashSenha(users.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, "error.hash.password")
		return
	}

	users.Password = string(hash)

	err = uc.UserUsecase.CreateUser(&users)
	if err != nil {
		c.JSON(http.StatusBadRequest, "error.create")
		return
	}

	c.JSON(http.StatusOK, "create.success")
}

func (uc *UserController) SearchAllUsers(c *gin.Context) {

	users, err := uc.UserUsecase.SearchAllUsers()
	if err != nil {
		c.JSON(http.StatusBadRequest, "error.search")
		return
	}
	if users == nil {
		c.JSON(http.StatusBadRequest, "not.found")
		return
	}

	c.JSON(200, users)
}

func (uc *UserController) SearchUserById(c *gin.Context) {

	id := c.Params.ByName("id")

	user, err := uc.UserUsecase.SearchUserById(conversion.StringToUint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, "error.search")
		return
	}

	if user.ID == 0 {
		c.JSON(http.StatusNotFound, "not.found")
		return
	}

	c.JSON(http.StatusOK, user)
}

func (uc *UserController) UpdateUser(c *gin.Context) {
	var user *models.User

	id := c.Params.ByName("id")

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, "error.json")
		return
	}

	err := uc.UserUsecase.ExistUserWithId(conversion.StringToUint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, "not.exist")
		return
	}

	if user.Password != "" {
		hash, err := security.HashSenha(user.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, "error.hash.passaword")
			return
		}

		user.Password = string(hash)
	}

	err = uc.UserUsecase.UpdateUser(user, conversion.StringToUint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, "error.update")
		return
	}

	c.JSON(http.StatusOK, "update.success")

}

func (uc *UserController) DeleteUser(c *gin.Context) {
	id := c.Params.ByName("id")

	err := uc.UserUsecase.ExistUserWithId(conversion.StringToUint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, "error.json")
		return
	}

	err = uc.UserUsecase.DeleteUser(conversion.StringToUint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, "error.delete")
		return
	}

	c.JSON(http.StatusOK, "delete.success")
}
