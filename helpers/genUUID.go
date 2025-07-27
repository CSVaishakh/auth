package helpers

import (
	"github.com/google/uuid"
)

func GenUUID () (string) {
	id := uuid.New().String()[:6]
	return id
}