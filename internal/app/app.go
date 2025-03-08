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
	opensearchSrv := service.NewOpenSearchService(opensearchRepo)
	cdcConsumer := consumer.NewOpenSearchCDC(opensearchRepo, *opensearchSrv)

	logs.Info("Starting the application")
	consumer.StartKafka(cdcConsumer)
	logs.Info("Application started successfully")

	time.Sleep(10 * time.Minute)

	// config.InitConfig()

	// app := fiber.New(fiber.Config{
	// 	ErrorHandler: func(c *fiber.Ctx, err error) error {
	// 		var statusCode int
	// 		var message string

	// 		// Check if error is of type *fiber.Error
	// 		var appErr errs.AppError
	// 		if errors.As(err, &appErr) {
	// 			statusCode = appErr.Code
	// 			message = appErr.Message
	// 		} else {
	// 			// If not, return a generic 500 status code
	// 			logs.Error(fmt.Sprintf("Unexpected error: %v", err))
	// 			statusCode = fiber.StatusInternalServerError
	// 			message = "Internal Server Error"
	// 		}

	// 		return c.Status(statusCode).JSON(types.GlobalErrorHandlerResp{
	// 			Success: false,
	// 			Message: message,
	// 		})
	// 	},
	// })

	// logs.Info(fmt.Sprintf("Server is running on port: %v", viper.GetInt("app.port")))
	// logs.Info(fmt.Sprintf("Server is running on port: %v", os.Getenv("APP_PORT")))
	// err := app.Listen(fmt.Sprintf(":%v", os.Getenv("APP_PORT")))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// brokers := os.Getenv("KAFKA_BROKERS")
	// // group_id := os.Getenv("KAFKA_GROUP_ID")
	// group_id := "g1"
	// topic := os.Getenv("KAFKA_TOPIC")

	// if brokers == "" || group_id == "" || topic == "" {
	// 	logs.Error("Missing Kafka environment variables")
	// 	return
	// }

	// config := kafka.ReaderConfig{
	// 	Brokers:  []string{brokers},
	// 	GroupID:  group_id,
	// 	Topic:    topic,
	// 	MaxBytes: 10e6, // 10MB
	// }

	// reader := kafka.NewReader(config)

	// for {
	// 	msg, err := reader.ReadMessage(context.Background())

	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	logs.Info(fmt.Sprintf("Message received: %v", string(msg.Value)))
	// }

	// Kafka reader health check and message consumption with timeout
	// go func() {
	// 	for {
	// 		// Set a timeout for the message consumption process
	// 		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // 10 second timeout
	// 		defer cancel()

	// 		msg, err := reader.ReadMessage(ctx)
	// 		if err != nil {
	// 			// Handle timeout or error
	// 			if errors.Is(err, context.DeadlineExceeded) {
	// 				logs.Warn("Timeout waiting for Kafka message")
	// 			} else {
	// 				logs.Error(fmt.Sprintf("Error reading from Kafka: %v", err))
	// 			}
	// 			return // exit the loop if error occurs
	// 		}

	// 		logs.Info(fmt.Sprintf("Kafka connection successful, message received: %v", string(msg.Value)))

	// 		if err := cdcConsumer.ConsumeMessage(msg); err != nil {
	// 			logs.Error(fmt.Sprintf("Error consuming message: %v", err))
	// 		}
	// 	}
	// }()

	// // Allow the server to keep running (do not return)
	// select {}

	// for {
	// 	msg, err := reader.ReadMessage(context.Background())
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	logs.Info(fmt.Sprintf("Message received: %v", string(msg.Value)))

	// 	if err := cdcConsumer.ConsumeMessage(msg); err != nil {
	// 		logs.Error(fmt.Sprintf("Error consuming message: %v", err))
	// 	}
	// }
}
