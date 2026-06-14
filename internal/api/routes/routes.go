package routes

import (
	"middleware/internal/api/middlewares"
	"middleware/internal/domain/interfaces"

	"github.com/casbin/casbin/v3"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func RouteConfiguration(db *gorm.DB, router *gin.Engine, enforcer *casbin.Enforcer, reloadNotifier interfaces.CLPReloadNotifier) {

	secretKey := viper.GetString("ACCESS_TOKEN")

	router.Use(middlewares.AccessPermissions(secretKey, enforcer))

	publicRoute := router.Group("")
	LoginRouter(db, publicRoute)
	RefreshRouter(db, publicRoute)

	protectedEquipamentoRouter := router.Group("/clps")
	protectedEquipamentoRouter.Use(middlewares.JwtAuthMiddleware(secretKey))
	ClpRoute(db, protectedEquipamentoRouter, reloadNotifier)

	protectedUsuariosRouter := router.Group("/users")
	protectedUsuariosRouter.Use(middlewares.JwtAuthMiddleware(secretKey))
	UsuarioRouter(db, protectedUsuariosRouter)

	protectedTagRouter := router.Group("/tags")
	protectedTagRouter.Use(middlewares.JwtAuthMiddleware(secretKey))
	TagRoute(db, protectedTagRouter, reloadNotifier)

	protectedSwapRouter := router.Group("/swaps")
	protectedSwapRouter.Use(middlewares.JwtAuthMiddleware(secretKey))
	SwapRoute(db, protectedSwapRouter)

	protectedTypeClpRouter := router.Group("/type-clps")
	protectedTypeClpRouter.Use(middlewares.JwtAuthMiddleware(secretKey))
	TypeClpRoute(db, protectedTypeClpRouter)

	protectedTypeTagRouter := router.Group("/type-tags")
	protectedTypeTagRouter.Use(middlewares.JwtAuthMiddleware(secretKey))
	TypeTagRoute(db, protectedTypeTagRouter)

	protectedTypeOperationRouter := router.Group("/type-operations")
	protectedTypeOperationRouter.Use(middlewares.JwtAuthMiddleware(secretKey))
	TypeOperationRoute(db, protectedTypeOperationRouter)
}
