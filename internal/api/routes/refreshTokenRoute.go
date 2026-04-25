package routes

import (
	"middleware/internal/api/controllers"
	"middleware/internal/repository"
	"middleware/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RefreshRouter(db *gorm.DB, group *gin.RouterGroup) {
	ur := repository.NewUserRepository(db)
	rtc := &controllers.RefreshTokenController{
		RefreshTokenUsecase: usecase.NewRefreshTokenUsecase(ur),
	}
	group.POST("/refresh/", rtc.RefreshToken)
}
