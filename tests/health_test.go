package handlers

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"webapp/api/handlers"
	"webapp/db"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var mockDB db.Database

func TestDBConnection(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Using system environment variables.")
	}
	host := os.Getenv("TEST_DB_HOST")
	user := os.Getenv("TEST_DB_USER")
	password := os.Getenv("TEST_DB_PASSWORD")
	dbName := os.Getenv("TEST_DB_NAME")
	port := os.Getenv("TEST_DB_PORT")
	dsn := fmt.Sprintf("host=" + host + " user=" + user + " password=" + password + " dbname=" + dbName + " port=" + port + " sslmode=disable")
	mdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	assert.NoError(t, err, "Failed to connect to database")

	sqlDB, err := mdb.DB()
	assert.NoError(t, err, "Failed to retrieve database connection")

	err = sqlDB.Ping()
	assert.NoError(t, err, "Database connection is not alive")

	err = mdb.AutoMigrate(&db.HealthCounter{})
	assert.NoError(t, err, "Failed to run AutoMigrate for HealthCounter")

	mockDB = &db.GormDatabase{DB: mdb}
}

func TestHealthCheckHandler_Success(t *testing.T) {
	TestDBConnection(t)
	_ = mockDB.Create(&db.HealthCounter{})
	handler := handlers.HealthCheckHandler(mockDB)

	req, _ := http.NewRequest("GET", "/healthz", nil)
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestHealthCheckHandler_Failure(t *testing.T) {
	TestDBConnection(t)
	gormDB := mockDB.(*db.GormDatabase).DB
	gormDB.Exec("DROP TABLE health_counters;")
	handler := handlers.HealthCheckHandler(mockDB)

	req, _ := http.NewRequest("GET", "/healthz", nil)
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusServiceUnavailable, resp.Code)
}

func TestHealthCheckHandler_405Failure(t *testing.T) {
	TestDBConnection(t)
	handler := handlers.HealthCheckHandler(mockDB)

	req, _ := http.NewRequest("POST", "/healthz", nil)
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusMethodNotAllowed, resp.Code)
}

func TestHealthCheckHandler_ContentLen(t *testing.T) {
	TestDBConnection(t)
	handler := handlers.HealthCheckHandler(mockDB)
	req, _ := http.NewRequest("GET", "/healthz", strings.NewReader("abc"))
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestHealthCheckHandler_QueryParam(t *testing.T) {
	TestDBConnection(t)
	handler := handlers.HealthCheckHandler(mockDB)
	req, _ := http.NewRequest("GET", "/healthz?test", nil)
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}
