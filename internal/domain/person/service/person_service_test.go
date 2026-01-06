package person

import (
	"errors"
	"testing"
	"time"

	personDto "pessoas-api/internal/contract/person"
	person "pessoas-api/internal/domain/person/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type repositoryMock struct {
	mock.Mock
}

func (r *repositoryMock) Save(person *person.Person) (ID int, err error) {
	args := r.Called(person)
	return args.Int(0), args.Error(1)
}

func (r *repositoryMock) FindAll(page, pageSize int, sortBy, sortOrder string) ([]*person.Person, int64, error) {
	args := r.Called(page, pageSize, sortBy, sortOrder)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*person.Person), args.Get(1).(int64), args.Error(2)
}

func (r *repositoryMock) FindByCPF(cpf string) (*person.Person, error) {
	args := r.Called(cpf)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*person.Person), args.Error(1)
}

func TestPersonService_CreatePerson_Success(t *testing.T) {
	assert := assert.New(t)
	repoMock := new(repositoryMock)

	repoMock.On("Save", mock.MatchedBy(func(person *person.Person) bool {
		if person == nil {
			return false
		}

		now := time.Now()

		return person.Name == "Jane Doe" &&
			person.CPF == "22233344405" &&
			person.Email == "jane.doe@example.com" &&
			person.PhoneNumber == "81998765432" &&
			person.BirthDate.Equal(time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC)) &&
			!person.CreatedAt.IsZero() &&
			!person.UpdatedAt.IsZero() &&
			person.CreatedAt.Equal(person.UpdatedAt) &&
			person.CreatedAt.Before(now.Add(time.Second)) &&
			person.CreatedAt.After(now.Add(-time.Second))
	})).Return(1, nil)

	service := NewPersonService(repoMock)

	createPersonDto := personDto.NewPersonDTO{
		Name:        "Jane Doe",
		CPF:         "222.333.444-05",
		BirthDate:   time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		PhoneNumber: "81 99876-5432",
		Email:       "jane.doe@example.com",
	}

	id, err := service.CreatePerson(createPersonDto)

	assert.Equal(1, id)
	assert.NoError(err)
	repoMock.AssertExpectations(t)
}

func TestPersonService_CreatePerson_RepoError(t *testing.T) {
	assert := assert.New(t)
	repoMock := new(repositoryMock)

	repoMock.On("Save", mock.Anything).Return(0, errors.New("repo error"))

	service := NewPersonService(repoMock)

	createPersonDto := personDto.NewPersonDTO{
		Name:        "Jane Doe",
		CPF:         "222.333.444-05",
		BirthDate:   time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		PhoneNumber: "81 99876-5432",
		Email:       "jane.doe@example.com",
	}

	_, err := service.CreatePerson(createPersonDto)

	assert.Error(err)
	repoMock.AssertExpectations(t)
}

func TestPersonService_CreatePerson_InputName(t *testing.T) {
	assert := assert.New(t)
	repoMock := new(repositoryMock)

	service := NewPersonService(repoMock)

	createPersonDto := personDto.NewPersonDTO{
		Name:        "",
		CPF:         "222.333.444-05",
		BirthDate:   time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		PhoneNumber: "81 99876-5432",
		Email:       "jane.doe@example.com",
	}

	_, err := service.CreatePerson(createPersonDto)

	assert.Error(err)
	repoMock.AssertExpectations(t)
}

func TestPersonService_CreatePerson_InputCPF(t *testing.T) {
	assert := assert.New(t)
	repoMock := new(repositoryMock)

	service := NewPersonService(repoMock)

	createPersonDto := personDto.NewPersonDTO{
		Name:        "Jane Doe",
		CPF:         "",
		BirthDate:   time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		PhoneNumber: "81 99876-5432",
		Email:       "jane.doe@example.com",
	}

	_, err := service.CreatePerson(createPersonDto)

	assert.Error(err)
	repoMock.AssertExpectations(t)
}

func TestPersonService_CreatePerson_InputPhoneNumber(t *testing.T) {
	assert := assert.New(t)
	repoMock := new(repositoryMock)

	service := NewPersonService(repoMock)

	createPersonDto := personDto.NewPersonDTO{
		Name:        "Jane Doe",
		CPF:         "222.333.444-05",
		BirthDate:   time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		PhoneNumber: "",
		Email:       "jane.doe@example.com",
	}

	_, err := service.CreatePerson(createPersonDto)

	assert.Error(err)
	repoMock.AssertExpectations(t)
}

func TestPersonService_CreatePerson_InputEmail(t *testing.T) {
	assert := assert.New(t)
	repoMock := new(repositoryMock)

	service := NewPersonService(repoMock)

	createPersonDto := personDto.NewPersonDTO{
		Name:        "Jane Doe",
		CPF:         "222.333.444-05",
		BirthDate:   time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		PhoneNumber: "81 99876-5432",
		Email:       "",
	}

	_, err := service.CreatePerson(createPersonDto)

	assert.Error(err)
	repoMock.AssertExpectations(t)
}

func TestPersonService_CreatePerson_InvalidCPF(t *testing.T) {
	assert := assert.New(t)
	repoMock := new(repositoryMock)

	service := NewPersonService(repoMock)

	createPersonDto := personDto.NewPersonDTO{
		Name:        "Jane Doe",
		CPF:         "123.456.789-00",
		BirthDate:   time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		PhoneNumber: "81 99876-5432",
		Email:       "jane.doe@example.com",
	}

	_, err := service.CreatePerson(createPersonDto)

	assert.Error(err)
	repoMock.AssertExpectations(t)
}

func TestPersonService_CreatePerson_InvalidEmail(t *testing.T) {
	assert := assert.New(t)
	repoMock := new(repositoryMock)

	service := NewPersonService(repoMock)

	createPersonDto := personDto.NewPersonDTO{
		Name:        "Jane Doe",
		CPF:         "222.333.444-05",
		BirthDate:   time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		PhoneNumber: "81 99876-5432",
		Email:       "jane.doe@invalid",
	}

	_, err := service.CreatePerson(createPersonDto)

	assert.Error(err)
	repoMock.AssertExpectations(t)
}

func TestPersonService_CreatePerson_InvalidPhoneNumber(t *testing.T) {
	assert := assert.New(t)
	repoMock := new(repositoryMock)

	service := NewPersonService(repoMock)

	createPersonDto := personDto.NewPersonDTO{
		Name:        "Jane Doe",
		CPF:         "222.333.444-05",
		BirthDate:   time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		PhoneNumber: "12345",
		Email:       "jane.doe@example.com",
	}

	_, err := service.CreatePerson(createPersonDto)

	assert.Error(err)
	repoMock.AssertExpectations(t)
}
