package handlers

import (
	"fmt"
	"go-auth-app/helpers"
	"go-auth-app/types"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func SignIn(c *fiber.Ctx) error {

	client, db_err := helpers.InItClient()
	var data map[string]string
	var users []types.User
	var user types.User
	var token types.Token
	var storedHashs []types.Secret
	var storedHash types.Secret
	expiry := 4 * time.Hour

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if db_err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": db_err.Error(),
		})
	}

	fmt.Println(data["email"])

	query_err := client.DB.From("users").Select("*").Execute(&users)
	fmt.Println(query_err)
	fmt.Println(len(users))
	if query_err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": query_err.Error(),
		})
	}

	if len(users) == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "no valid user",
		})
	}
	for i:=0;i<len(users);i++ {
		if users[i].Email == data["email"]{
			user = users[i]
		}
	}
	fmt.Println(user)


	query_err = client.DB.From("secrets").Select("*").Eq("userid",user.UserId).Execute(&storedHashs)
	if query_err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": query_err.Error(),
		})
	}
	storedHash = storedHashs[0]
	validation_err := helpers.ValidatePassword(data["password"],storedHash.Password)

	if validation_err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": validation_err.Error(),
		})
	}

	refreshToken, token_id, gen_err := helpers.GenJWT(user.UserId, user.Role, expiry, "refresh")

	if gen_err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": gen_err.Error(),
		})
	}

	token.Token_id = token_id
	token.UserId = user.UserId
	token.Role = user.Role
	token.Type = "refresh"
	token.Exp = strconv.FormatInt(time.Now().Add(expiry).Unix(),10)
	token.Iat = strconv.FormatInt(time.Now().Unix(),10)
	token.Status = true

	st_err := client.DB.From("jwt_tokens").Insert(token).Execute(nil)
	if st_err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": st_err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"refresh_token": refreshToken,
		"token_type":    "Bearer",
		"lifetime":      14400,
	})
}
