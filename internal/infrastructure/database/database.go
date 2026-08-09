package database

import (
	"fmt"
	"middleware/internal/domain/constants"
	"middleware/internal/infrastructure/database/migrations"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func InitializeDatabase() *gorm.DB {
	pathDir := GetPathDir()

	db, err := gorm.Open(sqlite.Open(filepath.Join(pathDir, constants.DB_NAME)), &gorm.Config{})
	if err != nil {
		fmt.Printf("Nao foi possivel conectar com o banco de dados: %v", err)
	} else {
		fmt.Printf("Banco conectado com sucesso")
	}

	migrations.RunMigrations(db)

	migrations.InitializeBasicUser(db)
	migrations.InitializeTypesClp(db)
	migrations.InitializeSwaps(db)
	migrations.InitializeTypesOperation(db)
	migrations.InitializeTypesTag(db)
	migrations.InitializeBasicUserProfile(db)

	return db
}

func GetPathDir() string {
	workingDir, err := os.Getwd()
	if err == nil {
		if _, statErr := os.Stat(filepath.Join(workingDir, constants.DB_NAME)); statErr == nil {
			return workingDir
		}
		if _, statErr := os.Stat(filepath.Join(workingDir, "go.mod")); statErr == nil {
			return workingDir
		}
	}

	execPath, err := os.Executable()
	if err != nil {
		fmt.Printf("Erro ao determinar o caminho do executavel: %v", err)
		return workingDir
	}

	return filepath.Dir(execPath)
}
