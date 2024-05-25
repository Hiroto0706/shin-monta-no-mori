package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword returns the bcrypt hash of the password
func HashPassword(password string) (string, error) {
	HashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(HashedPassword), nil
}

// CheckPassword checks if the provider password is correct or not
func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// CheckPassword checks if the provider email is correct or not
func CheckEmail(gotEmail string, email string) error {
	if gotEmail == email {
		return nil
	} else {
		return fmt.Errorf("%v is not a valid email", gotEmail)
	}
}
