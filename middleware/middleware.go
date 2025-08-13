package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"

	"os"
	"strings"
)

func VerifyToken (c *fiber.Ctx) error {
	fmt.Println("verifying")
	load_err := godotenv.Load()
	if load_err != nil { 
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error" : load_err.Error(),
		}) 
	}
	key := os.Getenv("JWT_SECRET")
	authHeader := c.Get("Authorization")

	if authHeader == "" { 
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error" : "Authheader not found",
		}) 
	}

	tokenStr := strings.TrimPrefix(authHeader,"Bearer ") 
	token,err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(key),nil
	})

	if err != nil || !token.Valid { 
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error" : "Unauthorized",
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok { 
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error" : "invalid token claims",
		})
	}

	token_id, ok := claims["token_id"].(string)
	if !ok { 
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error" : "invalid token_id",
		})
	}

	c.Locals("token_id",token_id)

	userid, ok := claims["sub"].(string)
	if !ok { 
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error" : "invalid userid",
		})
	}
	c.Locals("userid",userid)

	return c.Next()
}