package utils

import (
	"go-auth-app/helpers"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func SignIn (c *fiber.Ctx)  (string,error){
	
	client, db_err := helpers.InItClient()
	cred,err := helpers.DecodeJSON(c)
	var user User
	var hasshedPass string

	if err != nil { return "",err }
	if db_err != nil { return "",db_err }

	query_err := client.DB.From("users").Select("*").Eq("email", cred.Email).Eq("name",cred.Username).Execute(&user)
	if query_err != nil {
		return "", query_err
	}

	query_err = client.DB.From("secrets").Select("password").Eq("userid",user.UserId).Execute(&hasshedPass)
	if query_err != nil {
		return "", query_err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hasshedPass), []byte(cred.Password))

	if err != nil { return "", err }

	token,gen_err := helpers.GenJWT(user.UserId,user.Role,4*time.Hour,"refresh")

	if gen_err != nil { return "", gen_err }

	return token, err
}