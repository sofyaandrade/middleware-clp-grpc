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
	USUARIO_NAO_ENCONTRADO = "usuario.nao.encontrado"
)

func (rtc *RefreshTokenController) RefreshToken(c *gin.Context) {

	var refreshToken = viper.GetString("ACCESS_TOKEN")

	var request models.RefreshTokenRequest

	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, "error.json")
		return
	}

	id, _, err := rtc.RefreshTokenUsecase.ExtractIdToken(request.RefreshToken, refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "error.extract.id.token")
		return
	}

	usuario, err := rtc.RefreshTokenUsecase.SearchUserById(conversion.StringToUint(id))
	if err != nil {
		c.JSON(http.StatusUnauthorized, "not.found")
		return
	}

	accessToken, err := rtc.RefreshTokenUsecase.CreateAcessToken(usuario)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "error.create.access.token")
		return
	}

	refreshToken, err = rtc.RefreshTokenUsecase.CreateRefreshToken(usuario)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "error.create.refresh.token")
		return
	}

	refreshTokenResponse := models.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	c.JSON(http.StatusOK, refreshTokenResponse)
}
