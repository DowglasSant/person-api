package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(username, email, password string) (int, error) {
	args := m.Called(username, email, password)
	return args.Int(0), args.Error(1)
}

func (m *MockAuthService) Login(username, password string) (string, error) {
	args := m.Called(username, password)
	return args.String(0), args.Error(1)
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestRegister_Success(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)
	router := setupTestRouter()

	router.POST("/register", handler.Register)

	mockService.On("Register", "testuser", "test@example.com", "password123").Return(1, nil)

	requestBody := map[string]string{
		"username": "testuser",
		"email":    "test@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, float64(1), response["id"])
	assert.Equal(t, "Operator registered successfully", response["message"])
	mockService.AssertExpectations(t)
}

func TestRegister_InvalidJSON(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)
	router := setupTestRouter()

	router.POST("/register", handler.Register)

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "invalid_request", response["error"])
	assert.Contains(t, response["message"], "Invalid request body")
}

func TestRegister_MissingRequiredFields(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)
	router := setupTestRouter()

	router.POST("/register", handler.Register)

	requestBody := map[string]string{
		"username": "testuser",
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "invalid_request", response["error"])
}

func TestRegister_UsernameAlreadyExists(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)
	router := setupTestRouter()

	router.POST("/register", handler.Register)

	mockService.On("Register", "existinguser", "test@example.com", "password123").
		Return(0, errors.New("username already exists"))

	requestBody := map[string]string{
		"username": "existinguser",
		"email":    "test@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "registration_error", response["error"])
	assert.Equal(t, "username already exists", response["message"])
	mockService.AssertExpectations(t)
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)
	router := setupTestRouter()

	router.POST("/register", handler.Register)

	mockService.On("Register", "testuser", "existing@example.com", "password123").
		Return(0, errors.New("email already exists"))

	requestBody := map[string]string{
		"username": "testuser",
		"email":    "existing@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "registration_error", response["error"])
	assert.Equal(t, "email already exists", response["message"])
	mockService.AssertExpectations(t)
}

func TestRegister_ValidationError(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)
	router := setupTestRouter()

	router.POST("/register", handler.Register)

	requestBody := map[string]string{
		"username": "ab",
		"email":    "test@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "invalid_request", response["error"])
	assert.Contains(t, response["message"], "Username")
}

func TestRegister_ServiceError(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)
	router := setupTestRouter()

	router.POST("/register", handler.Register)

	mockService.On("Register", "testuser", "test@example.com", "password123").
		Return(0, errors.New("database connection failed"))

	requestBody := map[string]string{
		"username": "testuser",
		"email":    "test@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "registration_error", response["error"])
	mockService.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)
	router := setupTestRouter()

	router.POST("/login", handler.Login)

	mockService.On("Login", "testuser", "password123").Return("mock.jwt.token", nil)

	requestBody := map[string]string{
		"username": "testuser",
		"password": "password123",
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "mock.jwt.token", response["token"])
	assert.Equal(t, "Login successful", response["message"])
	mockService.AssertExpectations(t)
}

func TestLogin_InvalidJSON(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)
	router := setupTestRouter()

	router.POST("/login", handler.Login)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "invalid_request", response["error"])
	assert.Contains(t, response["message"], "Invalid request body")
}

func TestLogin_MissingCredentials(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)
	router := setupTestRouter()

	router.POST("/login", handler.Login)

	requestBody := map[string]string{
		"username": "testuser",
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "invalid_request", response["error"])
}

func TestLogin_InvalidCredentials(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)
	router := setupTestRouter()

	router.POST("/login", handler.Login)

	mockService.On("Login", "testuser", "wrongpassword").
		Return("", errors.New("invalid credentials"))

	requestBody := map[string]string{
		"username": "testuser",
		"password": "wrongpassword",
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "authentication_error", response["error"])
	assert.Equal(t, "invalid credentials", response["message"])
	mockService.AssertExpectations(t)
}

func TestLogin_UserNotFound(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)
	router := setupTestRouter()

	router.POST("/login", handler.Login)

	mockService.On("Login", "nonexistent", "password123").
		Return("", errors.New("invalid credentials"))

	requestBody := map[string]string{
		"username": "nonexistent",
		"password": "password123",
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "authentication_error", response["error"])
	mockService.AssertExpectations(t)
}

func TestLogin_InactiveAccount(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)
	router := setupTestRouter()

	router.POST("/login", handler.Login)

	mockService.On("Login", "inactiveuser", "password123").
		Return("", errors.New("operator account is inactive"))

	requestBody := map[string]string{
		"username": "inactiveuser",
		"password": "password123",
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "authentication_error", response["error"])
	assert.Equal(t, "operator account is inactive", response["message"])
	mockService.AssertExpectations(t)
}

func TestLogin_ServiceError(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)
	router := setupTestRouter()

	router.POST("/login", handler.Login)

	mockService.On("Login", "testuser", "password123").
		Return("", errors.New("failed to generate authentication token"))

	requestBody := map[string]string{
		"username": "testuser",
		"password": "password123",
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "authentication_error", response["error"])
	mockService.AssertExpectations(t)
}
