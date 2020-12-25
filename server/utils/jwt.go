package utils

import (
	"github.com/dgrijalva/jwt-go"
)

// DefaultTokenSecret is the default key used when there is no env vars
const DefaultTokenSecret = "D8AB934ECD92D0008477DFFE679D06A89779DE4BDB0D436A0B1D16E4F5BEFA2B"

// GPTokenSecret is the secret token for authorization through JWT
var GPTokenSecret = []byte(GetVar("GP_TOKEN_SECRET", DefaultTokenSecret))

// ExtractClaims will take the claim out of an authorization header.
// It takes in the JWT tokenString to read the claims from.
// It returns a map of string to interface for the claims.
func ExtractClaims(tokenString string) map[string]interface{} {
	if tokenString == "" {
		return nil
	}
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return GPTokenSecret, nil
	})
	return token.Claims.(jwt.MapClaims)
}
