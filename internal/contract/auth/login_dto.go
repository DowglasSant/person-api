package contract

type LoginDTO struct {
	Username string `json:"username" example:"john.doe" binding:"required"`       // Operator username
	Password string `json:"password" example:"SecurePass123!" binding:"required"` // Operator password
}

type LoginResponseDTO struct {
	Token   string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."` // JWT authentication token
	Message string `json:"message" example:"Login successful"`                      // Success message
}
