package operator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestNewOperator_Success(t *testing.T) {
	username := "johndoe"
	email := "john@example.com"
	password := "password123"

	op, err := NewOperator(username, email, password)

	assert.NoError(t, err)
	assert.NotNil(t, op)
	assert.Equal(t, username, op.Username)
	assert.Equal(t, email, op.Email)
	assert.NotEqual(t, password, op.PasswordHash)
	assert.True(t, op.Active)
	assert.NotZero(t, op.CreatedAt)
	assert.NotZero(t, op.UpdatedAt)

	err = bcrypt.CompareHashAndPassword([]byte(op.PasswordHash), []byte(password))
	assert.NoError(t, err)
}

func TestNewOperator_EmptyUsername(t *testing.T) {
	op, err := NewOperator("", "john@example.com", "password123")

	assert.Error(t, err)
	assert.Nil(t, op)
	assert.Equal(t, "username is required", err.Error())
}

func TestNewOperator_UsernameTooShort(t *testing.T) {
	op, err := NewOperator("ab", "john@example.com", "password123")

	assert.Error(t, err)
	assert.Nil(t, op)
	assert.Equal(t, "username must be at least 3 characters long", err.Error())
}

func TestNewOperator_UsernameTooLong(t *testing.T) {
	longUsername := "a"
	for i := 0; i < 51; i++ {
		longUsername += "a"
	}

	op, err := NewOperator(longUsername, "john@example.com", "password123")

	assert.Error(t, err)
	assert.Nil(t, op)
	assert.Equal(t, "username must not exceed 50 characters", err.Error())
}

func TestNewOperator_EmptyEmail(t *testing.T) {
	op, err := NewOperator("johndoe", "", "password123")

	assert.Error(t, err)
	assert.Nil(t, op)
	assert.Equal(t, "email is required", err.Error())
}

func TestNewOperator_EmailTooLong(t *testing.T) {
	longEmail := ""
	for i := 0; i < 101; i++ {
		longEmail += "a"
	}

	op, err := NewOperator("johndoe", longEmail, "password123")

	assert.Error(t, err)
	assert.Nil(t, op)
	assert.Equal(t, "email must not exceed 100 characters", err.Error())
}

func TestNewOperator_EmptyPassword(t *testing.T) {
	op, err := NewOperator("johndoe", "john@example.com", "")

	assert.Error(t, err)
	assert.Nil(t, op)
	assert.Equal(t, "password is required", err.Error())
}

func TestNewOperator_PasswordTooShort(t *testing.T) {
	op, err := NewOperator("johndoe", "john@example.com", "short")

	assert.Error(t, err)
	assert.Nil(t, op)
	assert.Equal(t, "password must be at least 8 characters long", err.Error())
}

func TestNewOperator_PasswordTooLong(t *testing.T) {
	longPassword := ""
	for i := 0; i < 73; i++ {
		longPassword += "a"
	}

	op, err := NewOperator("johndoe", "john@example.com", longPassword)

	assert.Error(t, err)
	assert.Nil(t, op)
	assert.Equal(t, "password must not exceed 72 characters", err.Error())
}

func TestNewOperator_UsernameExactly3Characters(t *testing.T) {
	op, err := NewOperator("abc", "john@example.com", "password123")

	assert.NoError(t, err)
	assert.NotNil(t, op)
	assert.Equal(t, "abc", op.Username)
}

func TestNewOperator_UsernameExactly50Characters(t *testing.T) {
	username := ""
	for i := 0; i < 50; i++ {
		username += "a"
	}

	op, err := NewOperator(username, "john@example.com", "password123")

	assert.NoError(t, err)
	assert.NotNil(t, op)
	assert.Equal(t, username, op.Username)
}

func TestNewOperator_EmailExactly100Characters(t *testing.T) {
	email := ""
	for i := 0; i < 100; i++ {
		email += "a"
	}

	op, err := NewOperator("johndoe", email, "password123")

	assert.NoError(t, err)
	assert.NotNil(t, op)
	assert.Equal(t, email, op.Email)
}

func TestNewOperator_PasswordExactly8Characters(t *testing.T) {
	op, err := NewOperator("johndoe", "john@example.com", "12345678")

	assert.NoError(t, err)
	assert.NotNil(t, op)
}

func TestNewOperator_PasswordExactly72Characters(t *testing.T) {
	password := ""
	for i := 0; i < 72; i++ {
		password += "a"
	}

	op, err := NewOperator("johndoe", "john@example.com", password)

	assert.NoError(t, err)
	assert.NotNil(t, op)
}

func TestValidatePassword_Success(t *testing.T) {
	password := "password123"
	op, err := NewOperator("johndoe", "john@example.com", password)

	assert.NoError(t, err)
	assert.True(t, op.ValidatePassword(password))
}

func TestValidatePassword_WrongPassword(t *testing.T) {
	password := "password123"
	op, err := NewOperator("johndoe", "john@example.com", password)

	assert.NoError(t, err)
	assert.False(t, op.ValidatePassword("wrongpassword"))
}

func TestValidatePassword_EmptyPassword(t *testing.T) {
	password := "password123"
	op, err := NewOperator("johndoe", "john@example.com", password)

	assert.NoError(t, err)
	assert.False(t, op.ValidatePassword(""))
}

func TestValidatePassword_CaseSensitive(t *testing.T) {
	password := "Password123"
	op, err := NewOperator("johndoe", "john@example.com", password)

	assert.NoError(t, err)
	assert.True(t, op.ValidatePassword("Password123"))
	assert.False(t, op.ValidatePassword("password123"))
	assert.False(t, op.ValidatePassword("PASSWORD123"))
}

func TestUpdatePassword_Success(t *testing.T) {
	op, err := NewOperator("johndoe", "john@example.com", "oldpassword123")
	assert.NoError(t, err)

	oldUpdatedAt := op.UpdatedAt
	newPassword := "newpassword123"

	err = op.UpdatePassword(newPassword)

	assert.NoError(t, err)
	assert.True(t, op.ValidatePassword(newPassword))
	assert.False(t, op.ValidatePassword("oldpassword123"))
	assert.True(t, op.UpdatedAt.After(oldUpdatedAt) || op.UpdatedAt.Equal(oldUpdatedAt))
}

func TestUpdatePassword_TooShort(t *testing.T) {
	op, err := NewOperator("johndoe", "john@example.com", "password123")
	assert.NoError(t, err)

	err = op.UpdatePassword("short")

	assert.Error(t, err)
	assert.Equal(t, "password must be at least 8 characters long", err.Error())
	assert.True(t, op.ValidatePassword("password123"))
}

func TestUpdatePassword_EmptyPassword(t *testing.T) {
	op, err := NewOperator("johndoe", "john@example.com", "password123")
	assert.NoError(t, err)

	err = op.UpdatePassword("")

	assert.Error(t, err)
	assert.Equal(t, "password must be at least 8 characters long", err.Error())
	assert.True(t, op.ValidatePassword("password123"))
}

func TestUpdatePassword_ExactlyMinLength(t *testing.T) {
	op, err := NewOperator("johndoe", "john@example.com", "password123")
	assert.NoError(t, err)

	err = op.UpdatePassword("12345678")

	assert.NoError(t, err)
	assert.True(t, op.ValidatePassword("12345678"))
	assert.False(t, op.ValidatePassword("password123"))
}

func TestValidateOperator_AllFieldsValid(t *testing.T) {
	err := validateOperator("johndoe", "john@example.com", "password123")
	assert.NoError(t, err)
}

func TestHashPassword_Success(t *testing.T) {
	password := "password123"
	hashedPassword, err := hashPassword(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)
	assert.NotEqual(t, password, hashedPassword)

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	assert.NoError(t, err)
}

func TestHashPassword_DifferentPasswordsProduceDifferentHashes(t *testing.T) {
	password1 := "password123"
	password2 := "password456"

	hash1, err1 := hashPassword(password1)
	hash2, err2 := hashPassword(password2)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEqual(t, hash1, hash2)
}

func TestHashPassword_SamePasswordProducesDifferentHashes(t *testing.T) {
	password := "password123"

	hash1, err1 := hashPassword(password)
	hash2, err2 := hashPassword(password)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEqual(t, hash1, hash2)

	err := bcrypt.CompareHashAndPassword([]byte(hash1), []byte(password))
	assert.NoError(t, err)

	err = bcrypt.CompareHashAndPassword([]byte(hash2), []byte(password))
	assert.NoError(t, err)
}

func TestOperator_FieldTypes(t *testing.T) {
	op, err := NewOperator("johndoe", "john@example.com", "password123")

	assert.NoError(t, err)
	assert.IsType(t, 0, op.ID)
	assert.IsType(t, "", op.Username)
	assert.IsType(t, "", op.Email)
	assert.IsType(t, "", op.PasswordHash)
	assert.IsType(t, true, op.Active)
}

func TestOperator_ActiveByDefault(t *testing.T) {
	op, err := NewOperator("johndoe", "john@example.com", "password123")

	assert.NoError(t, err)
	assert.True(t, op.Active)
}

func TestOperator_CanDeactivate(t *testing.T) {
	op, err := NewOperator("johndoe", "john@example.com", "password123")

	assert.NoError(t, err)
	assert.True(t, op.Active)

	op.Active = false
	assert.False(t, op.Active)
}
