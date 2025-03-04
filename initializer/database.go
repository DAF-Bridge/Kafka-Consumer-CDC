package initializer

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/DAF-Bridge/cdc-service/config"
	"github.com/DAF-Bridge/cdc-service/pkg/logs"
	"github.com/opensearch-project/opensearch-go"
	// "github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

var ESClient *opensearch.Client

func ConncectToDB() {
	// Load the configuration file
	config.InitConfig()

	// Using Viper to load the configuration file
	// dsn := fmt.Sprintf("%v://%v:%v@%v:%v/%v?sslmode=disable&TimeZone=Asia/Bangkok",
	// 	viper.GetString("dbRaspberry.driver"),
	// 	viper.GetString("dbRaspberry.user"),
	// 	viper.GetString("dbRaspberry.password"),
	// 	viper.GetString("dbRaspberry.host"),
	// 	viper.GetInt("dbRaspberry.port"),
	// 	viper.GetString("dbRaspberry.database"),
	// )

	dsn := os.Getenv("DATABASE_URL")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	DB = db

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database instance: ", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatal("Failed to ping database: ", err)
	}

	logs.Info("Successfully connected to PostgreSQL!")
}

func ConnectToOpensearch() *opensearch.Client {
	config.InitConfig()

	// bonsaiURL := viper.GetString("opensearch.host")
	// username := viper.GetString("opensearch.username")
	// password := viper.GetString("opensearch.password")

	bonsaiURL := os.Getenv("OPENSEARCH_HOST")
	username := os.Getenv("OPENSEARCH_USERNAME")
	password := os.Getenv("OPENSEARCH_PASSWORD")

	cfg := opensearch.Config{
		Addresses: []string{bonsaiURL},
		Username:  username,
		Password:  password,
	}

	// Create the client
	var err error
	ESClient, err = opensearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the Elasticsearch client: %s", err)
	}

	// Test connection
	res, err := ESClient.Info()
	if err != nil {
		log.Fatalf("Error getting info from Elasticsearch: %s", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalf("Error closing the response body: %s", err)
		}
	}(res.Body)

	status := res.Status()
	logs.Info(fmt.Sprintf("Successfully connected to Opensearch!, %s", status))

	return ESClient
}
