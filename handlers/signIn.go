package handlers

import (
	"go-auth-app/helpers"
	"go-auth-app/types"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func SignIn(c *fiber.Ctx) error {

	client, db_err := helpers.InItClient()
	cred, err := helpers.DecodeJSON(c)
	var user types.User
	var hasshedPass string
	var token types.Token
	expiry := 4 * time.Hour

	if err != nil { 
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error" : err.Error(),
		})	
	}
	
	if db_err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": db_err.Error(),
		})
	}

	query_err := client.DB.From("users").Select("*").Eq("email", cred.Email).Eq("name", cred.Username).Execute(&user)
	if query_err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":query_err.Error(),
		})
	}

	query_err = client.DB.From("secrets").Select("password").Eq("userid", user.UserId).Execute(&hasshedPass)
	if query_err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":query_err.Error(),
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(hasshedPass), []byte(cred.Password))

	if err != nil { 
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error" : err.Error(),
		})	
	}

	refreshToken, token_id, gen_err := helpers.GenJWT(user.UserId, user.Role, expiry, "refresh")

	if gen_err != nil { 
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error" :gen_err.Error(),
		})	
	}

	token.Token_id = token_id
	token.UserId = user.UserId
	token.Role = user.Role
	token.Type = "refresh"
	token.Exp = string(rune(time.Now().Add(expiry).Unix()))
	token.Iat = string(rune(time.Now().Unix()))
	token.Status = true

	query_err = client.DB.From("jwt_tokens").Insert(token).Execute(nil)
	if query_err != nil { 
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error" : query_err.Error(),
		})	
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"refresh_token": refreshToken,
			"token_type":    "Bearer",
			"lifetime":      14400,
	})
}
