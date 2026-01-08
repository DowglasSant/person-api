package handler

import (
	"log"
	"net/http"

	authContract "pessoas-api/internal/contract/auth"
	"pessoas-api/internal/domain/operator/ports"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService ports.AuthService
}

func NewAuthHandler(authService ports.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register creates a new operator account
func (h *AuthHandler) Register(c *gin.Context) {
	var dto authContract.RegisterDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		log.Printf("[ERROR] Register - Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Invalid request body: " + err.Error(),
		})
		return
	}

	id, err := h.authService.Register(dto.Username, dto.Email, dto.Password)
	if err != nil {
		log.Printf("[ERROR] Register - Registration failed for username %s: %v", dto.Username, err)

		statusCode := http.StatusUnprocessableEntity
		if err.Error() == "username already exists" || err.Error() == "email already exists" {
			statusCode = http.StatusConflict
		}

		c.JSON(statusCode, gin.H{
			"error":   "registration_error",
			"message": err.Error(),
		})
		return
	}

	log.Printf("[SUCCESS] Register - Operator created with ID: %d, Username: %s", id, dto.Username)
	c.JSON(http.StatusCreated, gin.H{
		"id":      id,
		"message": "Operator registered successfully",
	})
}

// Login authenticates an operator and returns a JWT token
func (h *AuthHandler) Login(c *gin.Context) {
	var dto authContract.LoginDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		log.Printf("[ERROR] Login - Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Invalid request body: " + err.Error(),
		})
		return
	}

	token, err := h.authService.Login(dto.Username, dto.Password)
	if err != nil {
		log.Printf("[ERROR] Login - Authentication failed for username %s: %v", dto.Username, err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "authentication_error",
			"message": err.Error(),
		})
		return
	}

	log.Printf("[SUCCESS] Login - Operator authenticated: %s", dto.Username)
	c.JSON(http.StatusOK, gin.H{
		"token":   token,
		"message": "Login successful",
	})
}
