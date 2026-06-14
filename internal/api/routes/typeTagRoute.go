package routes

import (
	"middleware/internal/api/controllers"
	"middleware/internal/repository"
	"middleware/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TypeTagRoute(db *gorm.DB, route *gin.RouterGroup) {
	mr := repository.NewTypeTagRepository(db)
	mc := &controllers.TypeTagController{
		TypeTagUsecase: usecase.NewTypeTagUsecase(mr),
	}
	route.GET("/", mc.SearchAllTypeTags)
}
