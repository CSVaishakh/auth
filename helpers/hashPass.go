package helpers

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPass (Password string) (string ,error) {
	passByte := []byte(Password)
	HasshedBytes, err := bcrypt.GenerateFromPassword(passByte,bcrypt.DefaultCost)
	return string(HasshedBytes), err
}