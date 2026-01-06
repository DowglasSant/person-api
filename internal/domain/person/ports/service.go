package ports

import (
	contract "pessoas-api/internal/contract/person"
	person "pessoas-api/internal/domain/person/model"
)

type PersonService interface {
	CreatePerson(dto contract.NewPersonDTO) (ID int, err error)
	ListPersons(page, pageSize int, sort, order string) ([]*person.Person, int64, error)
	FindPersonByCPF(cpf string) (*person.Person, error)
}
