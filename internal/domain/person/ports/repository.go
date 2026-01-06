package ports

import person "pessoas-api/internal/domain/person/model"

// PersonRepository defines the contract for person data persistence operations.
// This is the secondary port (driven port) that the domain requires to be implemented by adapters.
type PersonRepository interface {
	Save(person *person.Person) (ID int, err error)
	FindAll(page, size int, sortBy, sortOrder string) ([]*person.Person, int64, error)
	FindByCPF(cpf string) (*person.Person, error)
}
