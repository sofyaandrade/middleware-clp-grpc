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
		c.JSON(http.StatusBadRequest, "Dados da requisição inválidos.")
		return
	}

	usuario, err := lc.LoginUsecase.SearchUserByEmail(login.Email)

	if err != nil {
		c.JSON(http.StatusNotFound, "Nenhum usuário encontrado com este e-mail.")
		return
	}

	err = security.VerificaSenha(usuario.Password, login.Password)

	if err != nil {
		c.JSON(http.StatusUnauthorized, "Senha incorreta.")
		return
	}

	accessToken, err := lc.LoginUsecase.CreateAcessToken(usuario)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Não foi possível gerar o token de acesso.")
		return
	}

	refreshToken, err := lc.LoginUsecase.CreateRefreshToken(usuario)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Não foi possível gerar o token de renovação.")
		return
	}

	respostaLogin := models.ResponseLogin{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	c.JSON(http.StatusOK, respostaLogin)
}
