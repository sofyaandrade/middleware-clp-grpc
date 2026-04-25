package app

import (
	"log"
	"middleware/internal/infrastructure/database"
)

func InitializeProject() {

	db := database.InitializeDatabase()
	if db == nil {
		log.Fatal("Erro ao inicializar o banco de dados")
	}

	select {}
}
