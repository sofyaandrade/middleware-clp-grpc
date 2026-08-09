package routes

import (
	"middleware/internal/api/controllers"
	"middleware/internal/repository"
	"middleware/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserRoute(db *gorm.DB, group *gin.RouterGroup) {
	ur := repository.NewUserRepository(db)
	uc := &controllers.UserController{
		UserUsecase: usecase.NewUserUsecase(ur),
	}
	group.POST("/", uc.NewUser)
	group.GET("/", uc.SearchAllUsers)
	group.GET("/:id/", uc.SearchUserById)
	group.PATCH("/:id/", uc.UpdateUser)
	group.DELETE("/:id/", uc.DeleteUser)
}
