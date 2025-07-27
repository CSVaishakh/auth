package utils

import (
	"errors"
	"go-auth-app/helpers"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func SignOut(c *fiber.Ctx) (string,error) {
	load_err := godotenv.Load()
	client, err := helpers.InItClient()
	var refreshToken Token

	if load_err != nil { return "env error",load_err }
	if err != nil { return "",err }

	key := os.Getenv("JWT_SECRET")

	
	authHeader := c.Get("Authorization")

	if authHeader == "" { return "", errors.New("authorization header not found") }

	tokenStr := strings.TrimPrefix(authHeader,"Bearer ") 
	token,err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(key),nil
	})

	if err != nil || !token.Valid { return "",errors.New("Unauthorized")}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok { return "",errors.New("invalid claims") }
	tokenId:= claims["token_id"].(string)

	query_err := client.DB.From("jwt_tokens").Select("*").Eq("token_id",tokenId).Execute(&refreshToken)
	if query_err != nil { return  "", errors.New("Token retrival error")}

	return "SignOut Complete",nil
}