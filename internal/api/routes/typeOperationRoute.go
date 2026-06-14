package routes

import (
	"middleware/internal/api/controllers"
	"middleware/internal/repository"
	"middleware/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TypeOperationRoute(db *gorm.DB, route *gin.RouterGroup) {
	mr := repository.NewOperationTypeRepository(db)
	mc := &controllers.TypeOperationController{
		TypeOperationUsecase: usecase.NewOperationTypeUsecase(mr),
	}
	route.GET("/", mc.SearchAllTypeOperations)
}
