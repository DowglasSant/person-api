package contract

import "time"

// UpdatePersonDTO represents the data required to update a person
type UpdatePersonDTO struct {
	Name        string    `json:"name" example:"João Silva" binding:"required"`
	CPF         string    `json:"cpf" example:"111.444.777-35" binding:"required"`
	BirthDate   time.Time `json:"birth_date" example:"1990-01-15T00:00:00Z" binding:"required"`
	PhoneNumber string    `json:"phone" example:"81912345678" binding:"required"`
	Email       string    `json:"email" example:"joao.silva@email.com" binding:"required,email"`
}
