package person

import (
	personErr "pessoas-api/internal/domain/person/error"
	utils "pessoas-api/internal/domain/person/utils"
	"regexp"
	"strings"
	"time"
)

type Person struct {
	ID          int
	Name        string
	CPF         string
	BirthDate   time.Time
	PhoneNumber string
	Email       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewPerson(name string, cpf string, birthDate time.Time, phoneNumber string, email string) (*Person, error) {

	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)
	cpf = utils.OnlyDigits(cpf)
	phoneNumber = utils.OnlyDigits(phoneNumber)

	if name == "" {
		return nil, personErr.ErrNameRequired
	}
	if cpf == "" {
		return nil, personErr.ErrCPFRequired
	}
	if phoneNumber == "" {
		return nil, personErr.ErrPhoneRequired
	}
	if email == "" {
		return nil, personErr.ErrEmailRequired
	}

	if birthDate.IsZero() || birthDate.After(time.Now()) {
		return nil, personErr.ErrBirthDateInvalid
	}
	if !validateCPF(cpf) {
		return nil, personErr.ErrCPFInvalid
	}
	if !validateEmail(email) {
		return nil, personErr.ErrEmailInvalid
	}
	if !validatePhone(phoneNumber) {
		return nil, personErr.ErrPhoneInvalid
	}

	now := time.Now()

	return &Person{
		Name:        name,
		CPF:         cpf,
		BirthDate:   birthDate,
		PhoneNumber: phoneNumber,
		Email:       email,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func validateCPF(cpf string) bool {
	if len(cpf) != 11 {
		return false
	}

	allEqual := true
	for i := 1; i < 11; i++ {
		if cpf[i] != cpf[0] {
			allEqual = false
			break
		}
	}
	if allEqual {
		return false
	}

	if !checkCPFVerifier(cpf, 9, 10) {
		return false
	}

	if !checkCPFVerifier(cpf, 10, 11) {
		return false
	}

	return true
}

func checkCPFVerifier(cpf string, length int, weight int) bool {
	sum := 0

	for i := 0; i < length; i++ {
		sum += int(cpf[i]-'0') * (weight - i)
	}

	remainder := sum % 11
	var verifier byte

	if remainder < 2 {
		verifier = '0'
	} else {
		verifier = byte('0' + (11 - remainder))
	}

	return cpf[length] == verifier
}

func validateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	return emailRegex.MatchString(email)
}

func validatePhone(phone string) bool {
	if len(phone) != 10 && len(phone) != 11 {
		return false
	}

	return true
}
