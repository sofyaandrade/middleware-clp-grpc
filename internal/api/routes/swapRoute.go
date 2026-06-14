package routes

import (
	"middleware/internal/api/controllers"
	"middleware/internal/repository"
	"middleware/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SwapRoute(db *gorm.DB, route *gin.RouterGroup) {
	sr := repository.NewSwapRepository(db)
	mc := &controllers.SwapController{
		SwapUsecase: usecase.NewSwapUsecase(sr),
	}
	route.GET("/", mc.SearchAllSwaps)

}
