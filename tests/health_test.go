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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) Create(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

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
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	assert.NoError(t, err, "Failed to connect to database")

	sqlDB, err := db.DB()
	assert.NoError(t, err, "Failed to retrieve database connection")

	err = sqlDB.Ping()
	assert.NoError(t, err, "Database connection is not alive")
}

func TestHealthCheckHandler_Success(t *testing.T) {
	t.Parallel()
	mockDB := new(MockDatabase)
	mockDB.On("Create", mock.AnythingOfType("*db.HealthCounter")).Return(&gorm.DB{})
	handler := handlers.HealthCheckHandler(mockDB)

	req, _ := http.NewRequest("GET", "/healthz", nil)
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
	mockDB.AssertExpectations(t)
}

func TestHealthCheckHandler_Failure(t *testing.T) {
	t.Parallel()
	mockDB := new(MockDatabase)
	mockDB.On("Create", mock.AnythingOfType("*db.HealthCounter")).Return(&gorm.DB{Error: assert.AnError})
	handler := handlers.HealthCheckHandler(mockDB)

	req, _ := http.NewRequest("GET", "/healthz", nil)
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusServiceUnavailable, resp.Code)
	mockDB.AssertExpectations(t)
}

func TestHealthCheckHandler_405Failure(t *testing.T) {
	t.Parallel()
	mockDB := new(MockDatabase)
	handler := handlers.HealthCheckHandler(mockDB)

	req, _ := http.NewRequest("POST", "/healthz", nil)
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusMethodNotAllowed, resp.Code)
}

func TestHealthCheckHandler_ContentLen(t *testing.T) {
	t.Parallel()
	mockDB := new(MockDatabase)
	handler := handlers.HealthCheckHandler(mockDB)
	req, _ := http.NewRequest("GET", "/healthz", strings.NewReader("abc"))
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestHealthCheckHandler_QueryParam(t *testing.T) {
	t.Parallel()
	mockDB := new(MockDatabase)
	handler := handlers.HealthCheckHandler(mockDB)
	req, _ := http.NewRequest("GET", "/healthz?test", nil)
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}
