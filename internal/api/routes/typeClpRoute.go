package routes

import (
	"middleware/internal/api/controllers"
	"middleware/internal/repository"
	"middleware/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TypeClpRoute(db *gorm.DB, route *gin.RouterGroup) {
	mr := repository.NewTypeClpRepository(db)
	mc := &controllers.TypeClpController{
		TypeClpUsecase: usecase.NewTypeClpUsecase(mr),
	}
	route.GET("/", mc.SearchAllTypeClps)
}
