package middleware

import (
	"net/http"

	"github.com/UN0wen/pricewatch-vn/server/utils"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
)

// Authenticate checks for JWT token that is valid
func Authenticate(next http.Handler) http.Handler {
	JWTMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return utils.GPTokenSecret, nil
		},
		SigningMethod: jwt.SigningMethodHS256,
		UserProperty:  "auth",
	})
	return JWTMiddleware.Handler(next)
}
