package handlers

import (
	"net/http"
	"net/http/httptest"
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
	mockDB := new(MockDatabase)
	mockDB.On("Create", mock.AnythingOfType("*db.HealthCounter")).Return(&gorm.DB{})
	handler := handlers.HealthCheckHandler(mockDB)

	req, _ := http.NewRequest("GET", "/healthz", nil)
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
	mockDB.AssertExpectations(t)
}
