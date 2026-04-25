package database

import (
	"fmt"
	"log"
	"middleware/internal/domain/constants"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func InitializeDatabase() *gorm.DB {
	pathDir := GetPathDir()

	db, err := gorm.Open(sqlite.Open(filepath.Join(pathDir, constants.DB_NAME)), &gorm.Config{})
	if err != nil {
		fmt.Printf("Não foi possível conectar com o banco de dados: %v", err)
	} else {
		fmt.Printf("Banco conectado com sucesso")
	}

	return db
}

func GetPathDir() string {
	execPath, err := os.Executable()
	if err != nil {
		log.Fatalf("Erro ao determinar o caminho do executável: %v", err)
	}
	execDir := filepath.Dir(execPath)

	return filepath.Join(execDir, "")
}
