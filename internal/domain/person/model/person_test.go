package person

import (
	personErr "pessoas-api/internal/domain/person/error"
	utils "pessoas-api/internal/domain/person/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func validPersonInput() (name, cpf string, birthDate time.Time, phone, email string) {
	return "John Doe",
		"111.444.777-35",
		time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		"81 91234-5678",
		"john.doe@example.com"
}

func TestNewPerson_ShouldCreatePerson_WhenInputIsValid(t *testing.T) {
	assert := assert.New(t)

	name, cpf, birthDate, phone, email := validPersonInput()

	person, err := NewPerson(name, cpf, birthDate, phone, email)

	assert.NoError(err)
	assert.NotNil(person)

	assert.Equal(name, person.Name)
	assert.Equal(utils.OnlyDigits(cpf), person.CPF)
	assert.Equal(birthDate, person.BirthDate)
	assert.Equal(utils.OnlyDigits(phone), person.PhoneNumber)
	assert.Equal(email, person.Email)

	assert.WithinDuration(time.Now(), person.CreatedAt, time.Second)
	assert.WithinDuration(time.Now(), person.UpdatedAt, time.Second)
}

func TestNewPerson_ShouldFail_WhenCPFIsInvalid(t *testing.T) {
	assert := assert.New(t)

	name, _, birthDate, phone, email := validPersonInput()
	invalidCPF := "123.456.789-00"

	person, err := NewPerson(name, invalidCPF, birthDate, phone, email)

	assert.ErrorIs(err, personErr.ErrCPFInvalid)
	assert.Nil(person)
}

func TestNewPerson_ShouldFail_WhenEmailIsInvalid(t *testing.T) {
	assert := assert.New(t)

	name, cpf, birthDate, phone, _ := validPersonInput()
	invalidEmail := "john.doe@invalid"

	person, err := NewPerson(name, cpf, birthDate, phone, invalidEmail)

	assert.ErrorIs(err, personErr.ErrEmailInvalid)
	assert.Nil(person)
}

func TestNewPerson_ShouldFail_WhenPhoneNumberIsInvalid(t *testing.T) {
	assert := assert.New(t)

	name, cpf, birthDate, _, email := validPersonInput()
	invalidPhone := "12345"

	person, err := NewPerson(name, cpf, birthDate, invalidPhone, email)

	assert.ErrorIs(err, personErr.ErrPhoneInvalid)
	assert.Nil(person)
}

func TestNewPerson_ShouldFail_WhenNameIsEmpty(t *testing.T) {
	assert := assert.New(t)

	_, cpf, birthDate, phone, email := validPersonInput()
	emptyName := "   "

	person, err := NewPerson(emptyName, cpf, birthDate, phone, email)

	assert.ErrorIs(err, personErr.ErrNameRequired)
	assert.Nil(person)
}

func TestNewPerson_ShouldFail_WhenCPFIsEmpty(t *testing.T) {
	assert := assert.New(t)

	name, _, birthDate, phone, email := validPersonInput()
	emptyCPF := "   "

	person, err := NewPerson(name, emptyCPF, birthDate, phone, email)

	assert.ErrorIs(err, personErr.ErrCPFRequired)
	assert.Nil(person)
}

func TestNewPerson_ShouldFail_WhenPhoneNumberIsEmpty(t *testing.T) {
	assert := assert.New(t)

	name, cpf, birthDate, _, email := validPersonInput()
	emptyPhone := "   "

	person, err := NewPerson(name, cpf, birthDate, emptyPhone, email)

	assert.ErrorIs(err, personErr.ErrPhoneRequired)
	assert.Nil(person)
}

func TestNewPerson_ShouldFail_WhenEmailIsEmpty(t *testing.T) {
	assert := assert.New(t)

	name, cpf, birthDate, phone, _ := validPersonInput()
	emptyEmail := "   "

	person, err := NewPerson(name, cpf, birthDate, phone, emptyEmail)

	assert.ErrorIs(err, personErr.ErrEmailRequired)
	assert.Nil(person)
}

func TestNewPerson_ShouldFail_WhenBirthDateIsInFuture(t *testing.T) {
	assert := assert.New(t)

	name, cpf, _, phone, email := validPersonInput()
	futureBirthDate := time.Now().AddDate(1, 0, 0)

	person, err := NewPerson(name, cpf, futureBirthDate, phone, email)

	assert.ErrorIs(err, personErr.ErrBirthDateInvalid)
	assert.Nil(person)
}

func TestNewPerson_ShouldFail_WhenBirthDateIsZero(t *testing.T) {
	assert := assert.New(t)

	name, cpf, _, phone, email := validPersonInput()
	zeroBirthDate := time.Time{}

	person, err := NewPerson(name, cpf, zeroBirthDate, phone, email)

	assert.ErrorIs(err, personErr.ErrBirthDateInvalid)
	assert.Nil(person)
}
