package app

import (
	"os"
	"time"

	"github.com/DAF-Bridge/cdc-service/initializer"
	"github.com/DAF-Bridge/cdc-service/internal/consumer"
	"github.com/DAF-Bridge/cdc-service/internal/repository"
	"github.com/DAF-Bridge/cdc-service/internal/service"
	"github.com/DAF-Bridge/cdc-service/pkg/logs"
)

func init() {
	mode := os.Getenv("ENVIRONMENT")
	if mode != "production" && (mode == "" || mode == "dev") {
		initializer.LoadEnvVar()
	}
	initializer.ConncectToDB()
	initializer.ConnectToOpensearch()
	initializer.ConnectToMessageQueue()
}

func Start() {
	// Dependency injection
	// Initialize the repository
	opensearchRepo := repository.NewOpenSearchRepository(initializer.ESClient)
	jobRepo := repository.NewJobRepository(initializer.DB)
	eventRepo := repository.NewEventRepository(initializer.DB)
	orgRepo := repository.NewOrganizationRepository(initializer.DB)
	opensearchSrv := service.NewOpenSearchService(opensearchRepo, jobRepo, eventRepo, orgRepo)
	cdcConsumer := consumer.NewOpenSearchCDC(opensearchRepo, *opensearchSrv)

	logs.Info("Starting the application")
	consumer.StartKafka(cdcConsumer)
	logs.Info("Application started successfully")

	time.Sleep(10 * time.Minute)
}
