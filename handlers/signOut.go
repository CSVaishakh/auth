package utils

import (
	"errors"
	"go-auth-app/helpers"
	"go-auth-app/types"
	"github.com/gofiber/fiber/v2"

)

func SignOut(c *fiber.Ctx) (string,error) {
	client, err := helpers.InItClient()
	var refreshToken types.Token = 
	if err != nil { return err }


	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok { return "",errors.New("invalid claims") }
	tokenId:= claims["token_id"].(string)

	query_err := client.DB.From("jwt_tokens").Select("*").Eq("token_id",tokenId).Execute(&refreshToken)
	if query_err != nil { return  "", errors.New("token retrival error")}

	return "SignOut Complete",nil
}