package controllers

import (
	"middleware/internal/domain/interfaces"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserProfileController struct {
	UserProfileUsecase interfaces.UserProfileUsecase
}

func (uc *UserProfileController) SearchAllUserProfiles(c *gin.Context) {

	users, err := uc.UserProfileUsecase.SearchAllUserProfiles()
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
