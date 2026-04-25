package routes

import (
	"middleware/internal/api/middlewares"

	"github.com/casbin/casbin/v3"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func RouteConfiguration(db *gorm.DB, router *gin.Engine, enforcer *casbin.Enforcer) {

	secretKey := viper.GetString("ACCESS_TOKEN")

	router.Use(middlewares.AccessPermissions(secretKey, enforcer))

	publicRoute := router.Group("")
	LoginRouter(db, publicRoute)
	RefreshRouter(db, publicRoute)

	protectedEquipamentoRouter := router.Group("/clps")
	protectedEquipamentoRouter.Use(middlewares.JwtAuthMiddleware(secretKey))
	ClpRoute(db, protectedEquipamentoRouter)

	protectedUsuariosRouter := router.Group("/users")
	protectedUsuariosRouter.Use(middlewares.JwtAuthMiddleware(secretKey))
	UsuarioRouter(db, protectedUsuariosRouter)

	protectedTagRouter := router.Group("/tags")
	protectedTagRouter.Use(middlewares.JwtAuthMiddleware(secretKey))
	TagRoute(db, protectedTagRouter)
}
