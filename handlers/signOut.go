package handlers

import (
	"go-auth-app/utils"
	"github.com/gofiber/fiber/v2"

)

func SignOut(c *fiber.Ctx) error {
	client, err := utils.InItClient() 
	if err != nil { 
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error" : err.Error(),
		})	
	}

	tokenId := c.Locals("token_id").(string)
	userId := c.Locals("userid").(string) 

	query_err := client.DB.From("jwt_tokens").
	Update(map[string]interface{}{"status":false}).
	Eq("token_id",tokenId).
	Eq("userid",userId).
	Execute(nil)
	if query_err != nil { 
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error" : "error revoking token",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message" : "SignOut successful",
	})
}