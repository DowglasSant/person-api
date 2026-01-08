package service

import (
	"errors"
	"os"
	"testing"

	operator "pessoas-api/internal/domain/operator/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOperatorRepository struct {
	mock.Mock
}

func (m *MockOperatorRepository) Save(op *operator.Operator) (int, error) {
	args := m.Called(op)
	return args.Int(0), args.Error(1)
}

func (m *MockOperatorRepository) FindByUsername(username string) (*operator.Operator, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*operator.Operator), args.Error(1)
}

func (m *MockOperatorRepository) FindByEmail(email string) (*operator.Operator, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*operator.Operator), args.Error(1)
}

func (m *MockOperatorRepository) FindByID(id int) (*operator.Operator, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*operator.Operator), args.Error(1)
}

func TestRegister_Success(t *testing.T) {
	mockRepo := new(MockOperatorRepository)
	service := NewAuthService(mockRepo)

	mockRepo.On("FindByUsername", "newuser").Return(nil, nil)
	mockRepo.On("FindByEmail", "newuser@example.com").Return(nil, nil)
	mockRepo.On("Save", mock.AnythingOfType("*operator.Operator")).Return(1, nil)

	id, err := service.Register("newuser", "newuser@example.com", "password123")

	assert.NoError(t, err)
	assert.Equal(t, 1, id)
	mockRepo.AssertExpectations(t)
}

func TestRegister_UsernameAlreadyExists(t *testing.T) {
	mockRepo := new(MockOperatorRepository)
	service := NewAuthService(mockRepo)

	existingOp := &operator.Operator{
		ID:       1,
		Username: "existinguser",
		Email:    "existing@example.com",
		Active:   true,
	}

	mockRepo.On("FindByUsername", "existinguser").Return(existingOp, nil)

	id, err := service.Register("existinguser", "newuser@example.com", "password123")

	assert.Error(t, err)
	assert.Equal(t, 0, id)
	assert.Equal(t, "username already exists", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	mockRepo := new(MockOperatorRepository)
	service := NewAuthService(mockRepo)

	existingOp := &operator.Operator{
		ID:       1,
		Username: "existinguser",
		Email:    "existing@example.com",
		Active:   true,
	}

	mockRepo.On("FindByUsername", "newuser").Return(nil, nil)
	mockRepo.On("FindByEmail", "existing@example.com").Return(existingOp, nil)

	id, err := service.Register("newuser", "existing@example.com", "password123")

	assert.Error(t, err)
	assert.Equal(t, 0, id)
	assert.Equal(t, "email already exists", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestRegister_FindByUsernameError(t *testing.T) {
	mockRepo := new(MockOperatorRepository)
	service := NewAuthService(mockRepo)

	mockRepo.On("FindByUsername", "testuser").Return(nil, errors.New("database error"))

	id, err := service.Register("testuser", "test@example.com", "password123")

	assert.Error(t, err)
	assert.Equal(t, 0, id)
	assert.Equal(t, "failed to validate username", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestRegister_FindByEmailError(t *testing.T) {
	mockRepo := new(MockOperatorRepository)
	service := NewAuthService(mockRepo)

	mockRepo.On("FindByUsername", "testuser").Return(nil, nil)
	mockRepo.On("FindByEmail", "test@example.com").Return(nil, errors.New("database error"))

	id, err := service.Register("testuser", "test@example.com", "password123")

	assert.Error(t, err)
	assert.Equal(t, 0, id)
	assert.Equal(t, "failed to validate email", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestRegister_ValidationError_ShortUsername(t *testing.T) {
	mockRepo := new(MockOperatorRepository)
	service := NewAuthService(mockRepo)

	mockRepo.On("FindByUsername", "ab").Return(nil, nil)
	mockRepo.On("FindByEmail", "test@example.com").Return(nil, nil)

	id, err := service.Register("ab", "test@example.com", "password123")

	assert.Error(t, err)
	assert.Equal(t, 0, id)
	assert.Contains(t, err.Error(), "username must be at least 3 characters long")
	mockRepo.AssertExpectations(t)
}

func TestRegister_ValidationError_ShortPassword(t *testing.T) {
	mockRepo := new(MockOperatorRepository)
	service := NewAuthService(mockRepo)

	mockRepo.On("FindByUsername", "testuser").Return(nil, nil)
	mockRepo.On("FindByEmail", "test@example.com").Return(nil, nil)

	id, err := service.Register("testuser", "test@example.com", "short")

	assert.Error(t, err)
	assert.Equal(t, 0, id)
	assert.Contains(t, err.Error(), "password must be at least 8 characters long")
	mockRepo.AssertExpectations(t)
}

func TestRegister_SaveError(t *testing.T) {
	mockRepo := new(MockOperatorRepository)
	service := NewAuthService(mockRepo)

	mockRepo.On("FindByUsername", "newuser").Return(nil, nil)
	mockRepo.On("FindByEmail", "newuser@example.com").Return(nil, nil)
	mockRepo.On("Save", mock.AnythingOfType("*operator.Operator")).Return(0, errors.New("database error"))

	id, err := service.Register("newuser", "newuser@example.com", "password123")

	assert.Error(t, err)
	assert.Equal(t, 0, id)
	assert.Equal(t, "failed to create operator", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key-minimum-32-characters-long")
	defer os.Unsetenv("JWT_SECRET")

	mockRepo := new(MockOperatorRepository)
	service := NewAuthService(mockRepo)

	op, _ := operator.NewOperator("testuser", "test@example.com", "password123")
	op.ID = 1

	mockRepo.On("FindByUsername", "testuser").Return(op, nil)

	token, err := service.Login("testuser", "password123")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestLogin_UserNotFound(t *testing.T) {
	mockRepo := new(MockOperatorRepository)
	service := NewAuthService(mockRepo)

	mockRepo.On("FindByUsername", "nonexistent").Return(nil, nil)

	token, err := service.Login("nonexistent", "password123")

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, "invalid credentials", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestLogin_FindByUsernameError(t *testing.T) {
	mockRepo := new(MockOperatorRepository)
	service := NewAuthService(mockRepo)

	mockRepo.On("FindByUsername", "testuser").Return(nil, errors.New("database error"))

	token, err := service.Login("testuser", "password123")

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, "invalid credentials", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestLogin_InactiveOperator(t *testing.T) {
	mockRepo := new(MockOperatorRepository)
	service := NewAuthService(mockRepo)

	op, _ := operator.NewOperator("testuser", "test@example.com", "password123")
	op.ID = 1
	op.Active = false

	mockRepo.On("FindByUsername", "testuser").Return(op, nil)

	token, err := service.Login("testuser", "password123")

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, "operator account is inactive", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestLogin_InvalidPassword(t *testing.T) {
	mockRepo := new(MockOperatorRepository)
	service := NewAuthService(mockRepo)

	op, _ := operator.NewOperator("testuser", "test@example.com", "password123")
	op.ID = 1

	mockRepo.On("FindByUsername", "testuser").Return(op, nil)

	token, err := service.Login("testuser", "wrongpassword")

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, "invalid credentials", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestLogin_TokenGenerationError(t *testing.T) {
	os.Unsetenv("JWT_SECRET")

	mockRepo := new(MockOperatorRepository)
	service := NewAuthService(mockRepo)

	op, _ := operator.NewOperator("testuser", "test@example.com", "password123")
	op.ID = 1

	mockRepo.On("FindByUsername", "testuser").Return(op, nil)

	assert.Panics(t, func() {
		service.Login("testuser", "password123")
	})

	mockRepo.AssertExpectations(t)
}
