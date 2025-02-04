package handlers

import (
	"net/http"
	"net/http/httptest"
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
