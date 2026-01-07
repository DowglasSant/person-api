package contract

import "time"

// PersonResponseDTO represents person data returned by the API
type PersonResponseDTO struct {
	ID          int       `json:"id" example:"1"`                                       // Unique person ID
	Name        string    `json:"name" example:"João Silva"`                            // Full name
	CPF         string    `json:"cpf" example:"11144477735"`                            // Brazilian CPF (digits only)
	BirthDate   time.Time `json:"birth_date" example:"1990-01-15T00:00:00Z"`            // Date of birth
	PhoneNumber string    `json:"phone" example:"81912345678"`                          // Phone number (digits only)
	Email       string    `json:"email" example:"joao.silva@email.com"`                 // Email address
	CreatedAt   time.Time `json:"created_at" example:"2024-01-01T10:00:00Z"`            // Record creation timestamp
	UpdatedAt   time.Time `json:"updated_at" example:"2024-01-01T10:00:00Z"`            // Last update timestamp
}
