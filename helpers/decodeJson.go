package helpers

import (
		"go-auth-app/types"
	"github.com/gofiber/fiber/v2"
)

type Credentials = types.Credentials

func DecodeJSON (c *fiber.Ctx) (Credentials, error){
	var cred Credentials
	err := c.BodyParser(&cred)
	return  cred,err
}
