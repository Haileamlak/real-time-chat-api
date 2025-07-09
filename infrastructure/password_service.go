package infrastructure

import (
	"golang.org/x/crypto/bcrypt"
)

type PasswordService interface {
	HashPassword(password string) (string, error)
	ComparePasswords(hashedPassword string, password string) error
}

type passwordService struct{}

// creates a new password service
func NewPasswordService() PasswordService {
	return &passwordService{}
}

// hashes a password
func (s *passwordService) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// compares a hashed password with a plaintext password
func (s *passwordService) ComparePasswords(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}