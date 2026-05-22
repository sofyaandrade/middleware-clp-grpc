package routes

import (
	"middleware/internal/api/controllers"
	"middleware/internal/repository"
	"middleware/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ClpRoute(db *gorm.DB, route *gin.RouterGroup) {
	mr := repository.NewCLPRepository(db)
	mc := &controllers.CLPController{
		CLPUsecase: usecase.NewCLPUsecase(mr),
	}
	route.POST("/", mc.NewCLP)
	route.GET("/", mc.SearchAllClps)
	route.GET("/status/", mc.ClpsStatus)
	route.GET("/:id/", mc.SearchClpById)
	route.GET("/type/:idType/", mc.SearchClpByType)
	route.PATCH("/:id/", mc.UpdateClp)
	route.DELETE("/:id/", mc.DeleteClp)
}
