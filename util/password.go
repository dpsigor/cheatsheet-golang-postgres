package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword returns the bcrypt hash of the password
func HashPassword(pw string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %s", err)
	}
	return string(hashed), nil
}

// CheckPassword checks if the provided password is correct or not
func CheckPassword(pw string, hashed string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pw))
}
