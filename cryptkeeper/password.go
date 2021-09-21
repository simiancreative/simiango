package cryptkeeper

import (
	"golang.org/x/crypto/bcrypt"
)

func ComparePassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func EncryptPassword(text string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(text), 10)
	return string(bytes), err
}
