package handlers

import (
	"fmt"
	"go-auth-app/helpers"
	"go-auth-app/types"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func SignUp(c *fiber.Ctx) error {

	client, db_err := helpers.InItClient()
	var role_codes []types.RoleCode
	var data map[string]string
	var role string
	var user types.User
	var secret types.Secret

	if db_err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": db_err.Error(),
		})
	}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	query_err := client.DB.From("rolecodes").Select("*").Execute(&role_codes)
	if query_err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": query_err.Error(),
		})
	}

	for i := 0; i < len(role_codes); i++ {
		if data["role_code"] == role_codes[i].Code {
			role = role_codes[i].Role
			fmt.Println("Verified user role")
		}
	}

	user.Role = role
	user.UserId = helpers.GenUUID()
	user.Email = data["email"]
	user.Username = data["username"]
	user.CreatedAt = time.Now().Format(time.RFC3339)

	query_err = client.DB.From("users").Insert(user).Execute(nil)
	if query_err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": query_err.Error(),
		})
	}
	log.Println("Added data to user table")

	password := data["password"]
	hashedPass, err := helpers.HashPass(password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	secret.Password = hashedPass
	secret.UserId = user.UserId

	query_err = client.DB.From("secrets").Insert(secret).Execute(nil)
	if query_err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": query_err.Error(),
		})
	}
	log.Println("Added the secret")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "SignUp successful, Please Login",
	})
}
