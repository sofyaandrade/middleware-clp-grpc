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

	enforcer := AccessPermissionsConfig(db)
	if err := GinConfig(db, enforcer); err != nil {
		log.Fatalf("Erro ao iniciar servidor web: %v", err)
	}
}
