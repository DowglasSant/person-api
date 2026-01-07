package contract

import "time"

// NewPersonDTO represents the data required to create a new person
type NewPersonDTO struct {
	Name        string    `json:"name" example:"João Silva" binding:"required"`                             // Person's full name
	CPF         string    `json:"cpf" example:"111.444.777-35" binding:"required"`                          // Brazilian CPF (can be formatted or digits only)
	BirthDate   time.Time `json:"birth_date" example:"1990-01-15T00:00:00Z" binding:"required"`             // Date of birth
	PhoneNumber string    `json:"phone" example:"81912345678" binding:"required"`                           // Phone number (10 or 11 digits)
	Email       string    `json:"email" example:"joao.silva@email.com" binding:"required,email"`            // Valid email address
}
