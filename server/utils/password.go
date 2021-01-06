package utils

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a password string
func HashPassword(pass string) (hash []byte, err error) {
	hash, err = bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		err = errors.Wrapf(err, "Password hash failed")
		return
	}
	return
}
