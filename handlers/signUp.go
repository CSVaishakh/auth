package handlers

import (
	"fmt"
	"go-auth-app/utils"
	"go-auth-app/types"

	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func SignUp(c *fiber.Ctx) error {

	client, db_err := utils.InItClient()
	var data map[string]string
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


	user.UserId = utils.GenUUID()
	user.Email = data["email"]
	user.Username = data["name"]
	user.Role = ""
	user.CreatedAt = time.Now().Format(time.RFC3339)
	fmt.Println(user)

	query_err := client.DB.From("users").Insert(user).Execute(nil)
	if query_err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": query_err.Error(),
		})
	}
	log.Println("Added data to user table")

	password := data["password"]
	hashedPass, err := utils.HashPass(password)
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
		"message": "SignUp successful, Please Login after Role Assignmeent By Admin",
	})
}
