package setup

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var S3Client *s3.Client
var S3Bucket string

func DBConn() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Using system environment variables.")
	}
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	//dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + port + " connect_timeout=5"
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable connect_timeout=5", host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	log.Println("Connected to the database!")
	return db
}

func S3Conn() *s3.Client {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Using system environment variables.")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
	}

	S3Bucket = os.Getenv("S3_BUCKET")
	if S3Bucket == "" {
		log.Fatalf("S3_BUCKET is not set in the environment")
	}

	S3Client = s3.NewFromConfig(cfg)
	log.Println("Connected to S3!")
	return S3Client
}
