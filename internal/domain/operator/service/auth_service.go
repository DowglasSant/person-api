package service

import (
	"errors"
	"log"

	operator "pessoas-api/internal/domain/operator/model"
	"pessoas-api/internal/domain/operator/ports"
	"pessoas-api/internal/infrastructure/http/middleware"
)

type AuthServiceImpl struct {
	repository ports.OperatorRepository
}

func NewAuthService(repository ports.OperatorRepository) ports.AuthService {
	return &AuthServiceImpl{
		repository: repository,
	}
}

func (s *AuthServiceImpl) Register(username, email, password string) (int, error) {
	existingByUsername, err := s.repository.FindByUsername(username)
	if err != nil {
		log.Printf("[ERROR] Register - Failed to check username: %v", err)
		return 0, errors.New("failed to validate username")
	}
	if existingByUsername != nil {
		return 0, errors.New("username already exists")
	}

	existingByEmail, err := s.repository.FindByEmail(email)
	if err != nil {
		log.Printf("[ERROR] Register - Failed to check email: %v", err)
		return 0, errors.New("failed to validate email")
	}
	if existingByEmail != nil {
		return 0, errors.New("email already exists")
	}

	newOperator, err := operator.NewOperator(username, email, password)
	if err != nil {
		log.Printf("[ERROR] Register - Validation failed: %v", err)
		return 0, err
	}

	id, err := s.repository.Save(newOperator)
	if err != nil {
		log.Printf("[ERROR] Register - Failed to save operator: %v", err)
		return 0, errors.New("failed to create operator")
	}

	log.Printf("[SUCCESS] Register - Operator created with ID: %d, Username: %s", id, username)
	return id, nil
}

func (s *AuthServiceImpl) Login(username, password string) (string, error) {
	op, err := s.repository.FindByUsername(username)
	if err != nil {
		log.Printf("[ERROR] Login - Failed to find operator: %v", err)
		return "", errors.New("invalid credentials")
	}

	if op == nil {
		log.Printf("[WARN] Login - Operator not found: %s", username)
		return "", errors.New("invalid credentials")
	}

	if !op.Active {
		log.Printf("[WARN] Login - Inactive operator attempted login: %s", username)
		return "", errors.New("operator account is inactive")
	}

	if !op.ValidatePassword(password) {
		log.Printf("[WARN] Login - Invalid password for operator: %s", username)
		return "", errors.New("invalid credentials")
	}

	token, err := middleware.GenerateToken(op.ID, op.Username)
	if err != nil {
		log.Printf("[ERROR] Login - Failed to generate token: %v", err)
		return "", errors.New("failed to generate authentication token")
	}

	log.Printf("[SUCCESS] Login - Operator authenticated: %s (ID: %d)", username, op.ID)
	return token, nil
}
