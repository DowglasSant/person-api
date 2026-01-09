package person

import "errors"

var (
	ErrNameRequired     = errors.New("name is required")
	ErrCPFRequired      = errors.New("cpf is required")
	ErrCPFInvalid       = errors.New("cpf is invalid")
	ErrPhoneRequired    = errors.New("phone number is required")
	ErrPhoneInvalid     = errors.New("phone number is invalid")
	ErrEmailRequired    = errors.New("email is required")
	ErrEmailInvalid     = errors.New("email is invalid")
	ErrBirthDateInvalid = errors.New("birth date is invalid")
	ErrPersonNotFound   = errors.New("person not found")
)
