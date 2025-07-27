package main

import (
	"go-auth-app/utils"
	"log"
	"github.com/gofiber/fiber/v2"
)


func main () {
	app := fiber.New()

	app.Post("/signup",func (c *fiber.Ctx) error {

		msg, err := utils.SignUp(c)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
				"message": msg, 
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": msg,
		})
	})

	app.Post("/signin", func(c *fiber.Ctx) error {
		token, err := utils.SignIn(c)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error" : "Token not genrated",
				"details" : err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"refresh_token": token,
			"token_type": "Bearer",
			"lifetime": 14400, 
		})
	})

	log.Fatal(app.Listen(":5000"))
}