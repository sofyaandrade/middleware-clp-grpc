package controllers

import (
	"middleware/internal/domain/interfaces"
	"middleware/internal/domain/models"
	"middleware/internal/domain/security"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginController struct {
	LoginUsecase interfaces.LoginUsecase
}

func (lc *LoginController) Login(c *gin.Context) {
	var login models.Login
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, "erro.json")
		return
	}

	usuario, err := lc.LoginUsecase.SearchUserByEmail(login.Email)

	if err != nil {
		c.JSON(http.StatusNotFound, "not.foud.user.email")
		return
	}

	err = security.VerificaSenha(usuario.Password, login.Password)

	if err != nil {
		c.JSON(http.StatusUnauthorized, "error.password.verficate")
		return
	}

	accessToken, err := lc.LoginUsecase.CreateAcessToken(usuario)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "error.create.access.token")
		return
	}

	refreshToken, err := lc.LoginUsecase.CreateRefreshToken(usuario)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "erroe.create.refresh.token")
		return
	}

	respostaLogin := models.ResponseLogin{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	c.JSON(http.StatusOK, respostaLogin)
}
