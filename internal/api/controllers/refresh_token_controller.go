package controllers

import (
	"middleware/internal/domain/conversion"
	"middleware/internal/domain/interfaces"
	"middleware/internal/domain/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type RefreshTokenController struct {
	RefreshTokenUsecase interfaces.RefreshTokenUsecase
}

const (
	USUARIO_NAO_ENCONTRADO = "Usuário não encontrado."
)

func (rtc *RefreshTokenController) RefreshToken(c *gin.Context) {

	var refreshToken = viper.GetString("ACCESS_TOKEN")

	var request models.RefreshTokenRequest

	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Dados da requisição inválidos.")
		return
	}

	id, _, err := rtc.RefreshTokenUsecase.ExtractIdToken(request.RefreshToken, refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Token de renovação inválido ou expirado.")
		return
	}

	usuario, err := rtc.RefreshTokenUsecase.SearchUserById(conversion.StringToUint(id))
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Usuário não encontrado.")
		return
	}

	accessToken, err := rtc.RefreshTokenUsecase.CreateAcessToken(usuario)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Não foi possível gerar o token de acesso.")
		return
	}

	refreshToken, err = rtc.RefreshTokenUsecase.CreateRefreshToken(usuario)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Não foi possível gerar o token de renovação.")
		return
	}

	refreshTokenResponse := models.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	c.JSON(http.StatusOK, refreshTokenResponse)
}
