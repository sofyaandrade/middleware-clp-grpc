package routes

import (
	"middleware/internal/api/controllers"
	"middleware/internal/repository"
	"middleware/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func LoginRouter(db *gorm.DB, group *gin.RouterGroup) {
	ur := repository.NewUserRepository(db)
	lc := &controllers.LoginController{
		LoginUsecase: usecase.NewLoginUsecase(ur),
	}
	group.POST("/login/", lc.Login)
}
