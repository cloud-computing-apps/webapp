package setup

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var S3Client *s3.Client
var S3Bucket string

func DBConn() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Debug("No .env file found. Using system environment variables.")
	}
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + port + " connect_timeout=5"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	log.Info("Connected to the database!")
	return db
}

func S3Conn() *s3.Client {
	err := godotenv.Load()
	if err != nil {
		log.Debug("No .env file found. Using system environment variables.")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
	}

	S3Bucket = os.Getenv("S3_BUCKET")
	if S3Bucket == "" {
		log.Fatal("S3_BUCKET is not set in the environment")
	}

	S3Client = s3.NewFromConfig(cfg)
	log.Info("Connected to S3!")
	return S3Client
}
