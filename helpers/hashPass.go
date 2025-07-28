package helpers

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPass(Password string) (string, error) {
	passByte := []byte(Password)
	hashedBytes, err := bcrypt.GenerateFromPassword(passByte, bcrypt.DefaultCost)
	return string(hashedBytes), err
}

func ValidatePassword(password string, storedHash string) error {
    return bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
}