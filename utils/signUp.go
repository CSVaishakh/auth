package utils

import (
	"fmt"
	"go-auth-app/helpers"
	"go-auth-app/types"
	"time"
	"github.com/gofiber/fiber/v2"
	"log"
)

type Credentials = types.Credentials
type RoleCode = types.RoleCode
type User = types.User
type Secret = types.Secret

func SignUp(c *fiber.Ctx) (string, error) {

	client, db_err := helpers.InItClient()
	cred, err := helpers.DecodeJSON(c)
	
	var role_codes []RoleCode
	var user User
	var role string
	var secret Secret
	if err != nil {
		return " api request error", err
	}

	password := cred.Password
	hasshedPass, err := helpers.HashPass(password)

	if db_err != nil {
		return "db conecttion error", db_err
	}

	query_err := client.DB.From("rolecodes").Select("*").Execute(&role_codes)
	if query_err != nil {
		return "databse querying error", query_err
	}

	for i := 0; i < len(role_codes); i++ {
		if cred.Role_code == role_codes[i].Code {
			role = role_codes[i].Role
			fmt.Println("Verified user role")
		}
	}

	user.Role = role
	user.UserId = helpers.GenUUID()
	user.Email = cred.Email
	user.Username = cred.Username
	user.CreatedAt = time.Now().Format(time.RFC3339)

	query_err = client.DB.From("users").Insert(user).Execute(nil)
	if query_err != nil {
		return "Unable to create user, serverside error", query_err
	}
	log.Println("Added data to user table")

	secret.Password = hasshedPass
	secret.UserId = user.UserId

	query_err = client.DB.From("secrets").Insert(secret).Execute(nil)
	if query_err != nil {
		return "Unable to create user, serverside error", query_err
	}
	log.Println("Added the secret")


	return "User created successfully, please LogIn", err
}
