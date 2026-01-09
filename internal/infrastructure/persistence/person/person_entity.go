package person

import (
	"time"

	personModel "pessoas-api/internal/domain/person/model"
)

type PersonEntity struct {
	ID          int       `gorm:"column:id;primaryKey;autoIncrement"`
	Name        string    `gorm:"column:name;type:varchar(255);not null"`
	CPF         string    `gorm:"column:cpf;type:varchar(11);not null;uniqueIndex"`
	BirthDate   time.Time `gorm:"column:birth_date;type:date;not null"`
	PhoneNumber string    `gorm:"column:phone_number;type:varchar(11);not null"`
	Email       string    `gorm:"column:email;type:varchar(255);not null"`
	CreatedAt   time.Time `gorm:"column:created_at;type:timestamp;not null"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:timestamp;not null"`
}

func (PersonEntity) TableName() string {
	return "people.person"
}

func (e *PersonEntity) ToDomain() *personModel.Person {
	return &personModel.Person{
		ID:          e.ID,
		Name:        e.Name,
		CPF:         e.CPF,
		BirthDate:   e.BirthDate,
		PhoneNumber: e.PhoneNumber,
		Email:       e.Email,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func FromDomain(p *personModel.Person) *PersonEntity {
	return &PersonEntity{
		ID:          p.ID,
		Name:        p.Name,
		CPF:         p.CPF,
		BirthDate:   p.BirthDate,
		PhoneNumber: p.PhoneNumber,
		Email:       p.Email,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
