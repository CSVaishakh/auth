package main

import (
	"go-auth-app/handlers"
	"go-auth-app/middleware"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Post("/signout",handlers.SignUp)
	app.Post("/signout",handlers.SignIn)
	app.Post("/signout",middleware.VerifyToken,handlers.SignOut)
	
	log.Fatal(app.Listen(":5000"))
}
