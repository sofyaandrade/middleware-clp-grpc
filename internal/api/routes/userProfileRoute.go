package routes

import (
	"middleware/internal/api/controllers"
	"middleware/internal/repository"
	"middleware/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserProfileRoute(db *gorm.DB, group *gin.RouterGroup) {
	ur := repository.NewUserProfileRepository(db)
	uc := &controllers.UserProfileController{
		UserProfileUsecase: usecase.NewUserProfileUsecase(ur),
	}
	group.GET("/", uc.SearchAllUserProfiles)
}
