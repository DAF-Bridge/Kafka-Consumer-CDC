package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/DAF-Bridge/cdc-service/config"
	"github.com/DAF-Bridge/cdc-service/errs"
	"github.com/DAF-Bridge/cdc-service/initializer"
	"github.com/DAF-Bridge/cdc-service/internal/consumer"
	"github.com/DAF-Bridge/cdc-service/internal/repository"
	"github.com/DAF-Bridge/cdc-service/internal/service"
	"github.com/DAF-Bridge/cdc-service/pkg/logs"
	"github.com/DAF-Bridge/cdc-service/types"
	"github.com/gofiber/fiber/v2"
	"github.com/segmentio/kafka-go"
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
	config.InitConfig()

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			var statusCode int
			var message string

			// Check if error is of type *fiber.Error
			var appErr errs.AppError
			if errors.As(err, &appErr) {
				statusCode = appErr.Code
				message = appErr.Message
			} else {
				// If not, return a generic 500 status code
				logs.Error(fmt.Sprintf("Unexpected error: %v", err))
				statusCode = fiber.StatusInternalServerError
				message = "Internal Server Error"
			}

			return c.Status(statusCode).JSON(types.GlobalErrorHandlerResp{
				Success: false,
				Message: message,
			})
		},
	})

	// logs.Info(fmt.Sprintf("Server is running on port: %v", viper.GetInt("app.port")))
	logs.Info(fmt.Sprintf("Server is running on port: %v", os.Getenv("APP_PORT")))
	err := app.Listen(fmt.Sprintf(":%v", os.Getenv("APP_PORT")))
	if err != nil {
		log.Fatal(err)
	}

	// Dependency injection
	// Initialize the repository
	opensearchRepo := repository.NewOpenSearchRepository(initializer.ESClient)
	opensearchSrv := service.NewOpenSearchService(opensearchRepo)
	cdcConsumer := consumer.NewOpenSearchCDC(opensearchRepo, *opensearchSrv)

	brokers := os.Getenv("KAFKA_BROKERS")
	group_id := os.Getenv("KAFKA_GROUP_ID")
	topic := os.Getenv("KAFKA_TOPIC")

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{brokers},
		GroupID:  group_id,
		Topic:    topic,
		MinBytes: 10e2, // 1KB
		MaxBytes: 10e6, // 10MB
	})

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		logs.Info(fmt.Sprintf("Message received: %v", string(msg.Value)))

		if err := cdcConsumer.ConsumeMessage(msg); err != nil {
			logs.Error(fmt.Sprintf("Error consuming message: %v", err))
		}
	}
}
