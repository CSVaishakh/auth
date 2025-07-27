package helpers

import (
	"os"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func GenJWT (userId string,role string,expiry time.Duration, tokenType string) (string,string, error) {
	load_err := godotenv.Load()
	if load_err != nil { return "env error","",load_err }
	key := os.Getenv("JWT_SECRET")
	token_id := GenUUID()
	claims := jwt.MapClaims{
		"token_id": token_id,
		"sub":userId,
		"role":role,
		"type":tokenType,
		"exp": time.Now().Add(expiry).Unix(),
		"iat": time.Now().Unix(),
		"status": true,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	signedToken, err :=token.SignedString([]byte(key)) 

	if err != nil { return "","",err }

	return signedToken, token_id, err
}