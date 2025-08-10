package handlers

import (
	"go-auth-app/types"
	"go-auth-app/utils"

	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func AdminSignUp(c *fiber.Ctx) error {
	client, db_err := utils.InItClient()
	var data map[string]string
	var user types.Admin
	var secret types.Secret
	var licenses []struct {
		LicenseKey string `json:"license_key"`
	}
	var found bool

	if db_err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": db_err.Error(),
		})
	}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	user.UserId = utils.GenUUID()
	user.Email = data["email"]
	user.Username = data["name"]
	user.CreatedAt = time.Now().Format(time.RFC3339)
	user.Role = "Admin"
	user.LicenseKey = data["license_key"]

	query_err := client.DB.From("user_licenses").Select("license_key").Execute(&licenses)
	if query_err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": query_err.Error(),
		})
	}

	for _, l := range licenses {
		if user.LicenseKey == l.LicenseKey {
			found = true
			break
		}
	}

	if !found {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid license key",
		})
	}
	
	payload := map[string]interface{}{
		"userid":     user.UserId,
		"email":       user.Email,
		"name":    user.Username,
		"created_at":  user.CreatedAt,
		"role":        user.Role,
	}
	query_err = client.DB.From("users").Insert(payload).Execute(nil)
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
		"message": "SignUp successful, Please Login to setup you Organization",
	})
}