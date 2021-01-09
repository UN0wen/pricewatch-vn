package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
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

// GenerateJWT creates a JSON Web Token based on an id,
// with an expiration time of 1 day
// It returns the token string
func GenerateJWT(id uuid.UUID) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id.String(),
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	})
	tokenString, err := token.SignedString(GPTokenSecret)
	CheckError(err)
	return tokenString
}
