package main

import (
	"go-auth-app/handlers"
	"go-auth-app/middleware"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
	}))

	app.Post("/signup",handlers.SignUp)
	app.Post("/signin",handlers.SignIn)
	app.Post("/admin-signup",handlers.AdminSignUp)

	app.Post("/signout",middleware.VerifyToken,handlers.SignOut)
	app.Get("/profile",middleware.VerifyToken,handlers.GetProfile)

	app.Get("/verify",middleware.VerifyToken)	
	
	log.Fatal(app.Listen(":5000"))
}