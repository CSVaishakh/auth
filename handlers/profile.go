package handlers

import (
	"go-auth-app/types"

	"go-auth-app/utils"
	"github.com/gofiber/fiber/v2"
)

func GetProfile (c*fiber.Ctx) error {
	var user types.User
	client, err := utils.InItClient()
	if err != nil { 
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error" : err.Error(),
		})	
	}

	userId := c.Locals("userid").(string)

	query_err := client.DB.From("users").Select("*").Eq("userid",userId).Execute(&user)
	if query_err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": query_err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}