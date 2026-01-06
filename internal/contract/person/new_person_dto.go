package contract

import "time"

type NewPersonDTO struct {
	Name        string    `json:"name"`
	CPF         string    `json:"cpf"`
	BirthDate   time.Time `json:"birth_date"`
	PhoneNumber string    `json:"phone"`
	Email       string    `json:"email"`
}
