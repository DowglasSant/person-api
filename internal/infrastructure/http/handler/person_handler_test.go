package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	contract "pessoas-api/internal/contract/person"
	person "pessoas-api/internal/domain/person/model"
	"pessoas-api/internal/infrastructure/http/handler/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTest() (*gin.Engine, *mocks.MockPersonService) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockPersonService)
	handler := NewPersonHandler(mockService)

	router := gin.New()
	router.POST("/persons", handler.CreatePerson)
	router.GET("/persons", handler.ListPersons)
	router.GET("/persons/cpf/:cpf", handler.FindPersonByCPF)

	return router, mockService
}

// ========== CreatePerson Tests ==========

func TestCreatePerson_Success(t *testing.T) {
	router, mockService := setupTest()

	dto := contract.NewPersonDTO{
		Name:        "João Silva",
		CPF:         "111.444.777-35",
		BirthDate:   time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC),
		PhoneNumber: "81912345678",
		Email:       "joao.silva@email.com",
	}

	mockService.On("CreatePerson", dto).Return(1, nil)

	body, _ := json.Marshal(dto)
	req, _ := http.NewRequest("POST", "/persons", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, float64(1), response["id"])
	assert.Equal(t, "Person created successfully", response["message"])

	mockService.AssertExpectations(t)
}

func TestCreatePerson_InvalidJSON(t *testing.T) {
	router, mockService := setupTest()

	invalidJSON := []byte(`{"name": "João Silva", "cpf": }`) // Invalid JSON

	req, _ := http.NewRequest("POST", "/persons", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "invalid_request", response["error"])
	assert.Contains(t, response["message"], "Invalid request body")

	mockService.AssertNotCalled(t, "CreatePerson")
}

func TestCreatePerson_MissingRequiredFields(t *testing.T) {
	router, mockService := setupTest()

	dto := contract.NewPersonDTO{
		Name: "João Silva",
		// Missing required fields
	}

	body, _ := json.Marshal(dto)
	req, _ := http.NewRequest("POST", "/persons", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockService.AssertNotCalled(t, "CreatePerson")
}

func TestCreatePerson_ServiceValidationError(t *testing.T) {
	router, mockService := setupTest()

	dto := contract.NewPersonDTO{
		Name:        "João Silva",
		CPF:         "111.444.777-35",
		BirthDate:   time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC),
		PhoneNumber: "81912345678",
		Email:       "joao.silva@email.com",
	}

	mockService.On("CreatePerson", dto).Return(0, errors.New("invalid CPF"))

	body, _ := json.Marshal(dto)
	req, _ := http.NewRequest("POST", "/persons", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "validation_error", response["error"])
	assert.Equal(t, "invalid CPF", response["message"])

	mockService.AssertExpectations(t)
}

// ========== ListPersons Tests ==========

func TestListPersons_Success(t *testing.T) {
	router, mockService := setupTest()

	birthDate := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)
	createdAt := time.Now()

	persons := []*person.Person{
		{
			ID:          1,
			Name:        "João Silva",
			CPF:         "11144477735",
			BirthDate:   birthDate,
			PhoneNumber: "81912345678",
			Email:       "joao.silva@email.com",
			CreatedAt:   createdAt,
			UpdatedAt:   createdAt,
		},
		{
			ID:          2,
			Name:        "Maria Santos",
			CPF:         "22255588899",
			BirthDate:   birthDate,
			PhoneNumber: "81987654321",
			Email:       "maria.santos@email.com",
			CreatedAt:   createdAt,
			UpdatedAt:   createdAt,
		},
	}

	mockService.On("ListPersons", 1, 10, "id", "desc").Return(persons, int64(2), nil)

	req, _ := http.NewRequest("GET", "/persons?page=1&page_size=10&sort=id&order=desc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response contract.PaginatedResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, 1, response.Page)
	assert.Equal(t, 10, response.PageSize)
	assert.Equal(t, int64(2), response.TotalItems)
	assert.Equal(t, 1, response.TotalPages)

	mockService.AssertExpectations(t)
}

func TestListPersons_DefaultParameters(t *testing.T) {
	router, mockService := setupTest()

	mockService.On("ListPersons", 1, 10, "id", "desc").Return([]*person.Person{}, int64(0), nil)

	req, _ := http.NewRequest("GET", "/persons", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response contract.PaginatedResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, 1, response.Page)
	assert.Equal(t, 10, response.PageSize)
	assert.Equal(t, int64(0), response.TotalItems)
	assert.Equal(t, 0, response.TotalPages)

	mockService.AssertExpectations(t)
}

func TestListPersons_CustomPagination(t *testing.T) {
	router, mockService := setupTest()

	mockService.On("ListPersons", 2, 5, "name", "asc").Return([]*person.Person{}, int64(15), nil)

	req, _ := http.NewRequest("GET", "/persons?page=2&page_size=5&sort=name&order=asc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response contract.PaginatedResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, 2, response.Page)
	assert.Equal(t, 5, response.PageSize)
	assert.Equal(t, int64(15), response.TotalItems)
	assert.Equal(t, 3, response.TotalPages) // ceil(15/5) = 3

	mockService.AssertExpectations(t)
}

func TestListPersons_ServiceError(t *testing.T) {
	router, mockService := setupTest()

	mockService.On("ListPersons", 1, 10, "id", "desc").
		Return(nil, int64(0), errors.New("database connection error"))

	req, _ := http.NewRequest("GET", "/persons", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "internal_error", response["error"])
	assert.Contains(t, response["message"], "database connection error")

	mockService.AssertExpectations(t)
}

func TestListPersons_EmptyResult(t *testing.T) {
	router, mockService := setupTest()

	mockService.On("ListPersons", 1, 10, "id", "desc").Return([]*person.Person{}, int64(0), nil)

	req, _ := http.NewRequest("GET", "/persons", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response contract.PaginatedResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, int64(0), response.TotalItems)
	assert.Equal(t, 0, response.TotalPages)

	mockService.AssertExpectations(t)
}

// ========== FindPersonByCPF Tests ==========

func TestFindPersonByCPF_Success(t *testing.T) {
	router, mockService := setupTest()

	birthDate := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)
	createdAt := time.Now()

	personObj := &person.Person{
		ID:          1,
		Name:        "João Silva",
		CPF:         "11144477735",
		BirthDate:   birthDate,
		PhoneNumber: "81912345678",
		Email:       "joao.silva@email.com",
		CreatedAt:   createdAt,
		UpdatedAt:   createdAt,
	}

	mockService.On("FindPersonByCPF", "111.444.777-35").Return(personObj, nil)

	req, _ := http.NewRequest("GET", "/persons/cpf/111.444.777-35", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response person.Person
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, 1, response.ID)
	assert.Equal(t, "João Silva", response.Name)
	assert.Equal(t, "11144477735", response.CPF)

	mockService.AssertExpectations(t)
}

func TestFindPersonByCPF_NotFound(t *testing.T) {
	router, mockService := setupTest()

	mockService.On("FindPersonByCPF", "111.444.777-35").Return(nil, nil)

	req, _ := http.NewRequest("GET", "/persons/cpf/111.444.777-35", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "not_found", response["error"])
	assert.Equal(t, "Person not found with the provided CPF", response["message"])

	mockService.AssertExpectations(t)
}

func TestFindPersonByCPF_ServiceError(t *testing.T) {
	router, mockService := setupTest()

	mockService.On("FindPersonByCPF", "111.444.777-35").
		Return(nil, errors.New("database error"))

	req, _ := http.NewRequest("GET", "/persons/cpf/111.444.777-35", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "internal_error", response["error"])
	assert.Contains(t, response["message"], "database error")

	mockService.AssertExpectations(t)
}

func TestFindPersonByCPF_WithFormattedCPF(t *testing.T) {
	router, mockService := setupTest()

	personObj := &person.Person{
		ID:   1,
		Name: "João Silva",
		CPF:  "11144477735",
	}

	mockService.On("FindPersonByCPF", "111.444.777-35").Return(personObj, nil)

	req, _ := http.NewRequest("GET", "/persons/cpf/111.444.777-35", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestFindPersonByCPF_WithUnformattedCPF(t *testing.T) {
	router, mockService := setupTest()

	personObj := &person.Person{
		ID:   1,
		Name: "João Silva",
		CPF:  "11144477735",
	}

	mockService.On("FindPersonByCPF", "11144477735").Return(personObj, nil)

	req, _ := http.NewRequest("GET", "/persons/cpf/11144477735", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// ========== Edge Cases ==========

func TestNewPersonHandler(t *testing.T) {
	mockService := new(mocks.MockPersonService)
	handler := NewPersonHandler(mockService)

	assert.NotNil(t, handler)
	assert.Equal(t, mockService, handler.service)
}

func TestListPersons_CalculatesTotalPagesCorrectly(t *testing.T) {
	router, mockService := setupTest()

	testCases := []struct {
		name          string
		pageSize      int
		totalItems    int64
		expectedPages int
	}{
		{"exact division", 10, 50, 5},
		{"with remainder", 10, 55, 6},
		{"less than page size", 10, 5, 1},
		{"empty result", 10, 0, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService.On("ListPersons", 1, tc.pageSize, "id", "desc").
				Return([]*person.Person{}, tc.totalItems, nil).
				Once()

			req, _ := http.NewRequest("GET", "/persons?page_size="+strconv.Itoa(tc.pageSize), nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			var response contract.PaginatedResponse
			json.Unmarshal(w.Body.Bytes(), &response)

			assert.Equal(t, tc.expectedPages, response.TotalPages)
		})
	}
}
