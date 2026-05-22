package app

import (
	"fmt"
	"middleware/internal/api/middlewares"
	"middleware/internal/api/routes"
	"os"
	"path/filepath"

	"github.com/casbin/casbin/v3"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func GinConfig(db *gorm.DB, enforcer *casbin.Enforcer) error {
	serverPort := ":1710" //api port

	gin := gin.Default()
	gin.Use(middlewares.CORSMiddleware())
	routes.RouteConfiguration(db, gin, enforcer)

	gin.GET("/documentacao/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	if err := gin.Run(serverPort); err != nil {
		fmt.Println("")
		return err
	}
	return nil

}

func GinConfigFront() error {
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("")
		return err
	}

	exeDir := filepath.Dir(exePath)
	absPath := filepath.Join(exeDir, "build")

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		fmt.Println("")
		return err
	}

	frontEndPort := fmt.Sprintf(":%d", 2910) //api port
	ginFrontEnd := gin.Default()

	ginFrontEnd.Use(static.Serve("/", static.LocalFile(absPath, true)))
	if err := ginFrontEnd.Run(frontEndPort); err != nil {
		fmt.Println("")
		return err
	}

	return nil
}
