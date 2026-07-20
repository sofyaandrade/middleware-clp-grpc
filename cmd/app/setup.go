package app

import (
	"context"
	"fmt"
	"log"
	"middleware/internal/infrastructure/clp"
	"middleware/internal/infrastructure/database"
	modbusmaster "middleware/internal/infrastructure/modbusMaster"
	"middleware/internal/repository"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func InitializeProject() {

	db := database.InitializeDatabase()
	if db == nil {
		log.Fatal("Erro ao inicializar o banco de dados")
	}

	appCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var jobsWG sync.WaitGroup

	clpRepository := repository.NewCLPRepository(db)
	modbusMasterService := modbusmaster.NewService(clpRepository)
	clpManager := clp.NewManager(
		modbusMasterService,
	)
	clpManager.Start(appCtx, &jobsWG)

	enforcer := AccessPermissionsConfig(db)

	serverErrorChan := make(chan error, 3)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(quit)

	go func() {
		if err := GinConfigFront(); err != nil {
			serverErrorChan <- fmt.Errorf("erro ao iniciar front-end: %w", err)
		}
	}()

	go func() {
		if err := GinConfig(db, enforcer, modbusMasterService); err != nil {
			serverErrorChan <- fmt.Errorf("erro ao iniciar servidor web: %w", err)
		}
	}()

	go func() {
		if err := GRPCConfig(appCtx); err != nil {
			serverErrorChan <- fmt.Errorf("erro ao iniciar servidor gRPC: %w", err)
		}
	}()

	select {
	case err := <-serverErrorChan:
		log.Printf("Erro critico detectado: %v", err)
	case <-quit:
		log.Println("Shutdown recebido")
	}

	cancel()
	jobsWG.Wait()
}
