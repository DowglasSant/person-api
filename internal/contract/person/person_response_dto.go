package contract

import "time"

type PersonResponseDTO struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	CPF         string    `json:"cpf"`
	BirthDate   time.Time `json:"birth_date"`
	PhoneNumber string    `json:"phone"`
	Email       string    `json:"email"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
