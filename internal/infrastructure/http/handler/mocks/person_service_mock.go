package mocks

import (
	contract "pessoas-api/internal/contract/person"
	person "pessoas-api/internal/domain/person/model"

	"github.com/stretchr/testify/mock"
)

// MockPersonService is a mock implementation of ports.PersonService
type MockPersonService struct {
	mock.Mock
}

func (m *MockPersonService) CreatePerson(dto contract.NewPersonDTO) (int, error) {
	args := m.Called(dto)
	return args.Int(0), args.Error(1)
}

func (m *MockPersonService) ListPersons(page, pageSize int, sort, order string) ([]*person.Person, int64, error) {
	args := m.Called(page, pageSize, sort, order)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*person.Person), args.Get(1).(int64), args.Error(2)
}

func (m *MockPersonService) FindPersonByCPF(cpf string) (*person.Person, error) {
	args := m.Called(cpf)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*person.Person), args.Error(1)
}
