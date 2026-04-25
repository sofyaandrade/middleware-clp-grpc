package routes

import (
	"middleware/internal/api/controllers"
	"middleware/internal/repository"
	"middleware/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TagRoute(db *gorm.DB, route *gin.RouterGroup) {
	mr := repository.NewTagRepository(db)
	mc := &controllers.TagController{
		TagUsecase: usecase.NewTagUsecase(mr),
	}
	route.POST("/", mc.NewTag)
	route.GET("/", mc.SearchAllTags)
	route.GET("/:id/", mc.SearchTagById)
	route.PATCH("/:id/", mc.UpdateTag)
	route.DELETE("/:id/", mc.DeleteTag)
}
