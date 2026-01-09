package operator

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Operator struct {
	ID           int       `gorm:"primaryKey"`
	Username     string    `gorm:"uniqueIndex;not null"`
	Email        string    `gorm:"uniqueIndex;not null"`
	PasswordHash string    `gorm:"column:password_hash;not null"`
	Active       bool      `gorm:"default:true"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

func NewOperator(username, email, password string) (*Operator, error) {
	if err := validateOperator(username, email, password); err != nil {
		return nil, err
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	return &Operator{
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
		Active:       true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

func (o *Operator) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(o.PasswordHash), []byte(password))
	return err == nil
}

func (o *Operator) UpdatePassword(newPassword string) error {
	if len(newPassword) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	hashedPassword, err := hashPassword(newPassword)
	if err != nil {
		return err
	}

	o.PasswordHash = hashedPassword
	o.UpdatedAt = time.Now()
	return nil
}

func validateOperator(username, email, password string) error {
	if username == "" {
		return errors.New("username is required")
	}

	if len(username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}

	if len(username) > 50 {
		return errors.New("username must not exceed 50 characters")
	}

	if email == "" {
		return errors.New("email is required")
	}

	if len(email) > 100 {
		return errors.New("email must not exceed 100 characters")
	}

	if password == "" {
		return errors.New("password is required")
	}

	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	if len(password) > 72 {
		return errors.New("password must not exceed 72 characters")
	}

	return nil
}

func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("failed to hash password")
	}
	return string(hashedBytes), nil
}
