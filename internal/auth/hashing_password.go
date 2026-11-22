package auth

import (
	"fmt"
	"log"

	"github.com/alexedwards/argon2id"
)

func HashPassword(plainText string) (string, error) {
	hashed_password, err := argon2id.CreateHash(plainText, argon2id.DefaultParams)
	if err != nil {
		log.Println("The creation of a hash caused an error")
		return "", fmt.Errorf("Couldn't create the hash")
	}

	return hashed_password, nil
}

func CheckPassword(plainText, hashed_password string) (bool, error) {
	ok, err := argon2id.ComparePasswordAndHash(plainText, hashed_password)
	if err != nil {
		log.Println("Error was caused when checking the password")
		return false, fmt.Errorf("Couldn't verify the password")
	}
	return ok, nil
}
