package app

import (
	"fmt"
	"middleware/internal/api/middlewares"
	"middleware/internal/api/routes"
	"middleware/internal/domain/interfaces"
	"os"
	"path/filepath"

	"github.com/casbin/casbin/v3"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func GinConfig(db *gorm.DB, enforcer *casbin.Enforcer, reloadNotifier interfaces.CLPReloadNotifier) error {
	serverPort := ":1710" //api port

	gin := gin.Default()
	gin.Use(middlewares.CORSMiddleware())
	routes.RouteConfiguration(db, gin, enforcer, reloadNotifier)

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
	absPath, err := findFrontendBuildDir(exeDir)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	frontEndPort := fmt.Sprintf(":%d", 2910)
	ginFrontEnd := gin.Default()

	ginFrontEnd.Static("/assets", filepath.Join(absPath, "assets"))
	ginFrontEnd.GET("/", func(c *gin.Context) {
		c.File(filepath.Join(absPath, "index.html"))
	})
	ginFrontEnd.NoRoute(func(c *gin.Context) {
		requestedPath := filepath.Clean("." + c.Request.URL.Path)
		filePath := filepath.Join(absPath, requestedPath)

		if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
			c.File(filePath)
			return
		}

		c.File(filepath.Join(absPath, "index.html"))
	})

	if err := ginFrontEnd.Run(frontEndPort); err != nil {
		fmt.Println("")
		return err
	}

	return nil
}

func findFrontendBuildDir(exeDir string) (string, error) {
	candidates := []string{
		filepath.Join(exeDir, "dist"),
		filepath.Join(exeDir, "build"),
	}

	for _, candidate := range candidates {
		indexPath := filepath.Join(candidate, "index.html")
		if info, err := os.Stat(indexPath); err == nil && !info.IsDir() {
			return candidate, nil
		}
	}

	return "", fmt.Errorf(
		"Pasta do front-end nao encontrada. Coloque a pasta dist ao lado do executavel em %s. Iniciando apenas a API.",
		filepath.Join(exeDir, "dist"),
	)
}
