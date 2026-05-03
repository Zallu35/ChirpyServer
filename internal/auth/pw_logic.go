package auth

import (
	"fmt"

	"github.com/alexedwards/argon2id"
)

func HashPassword(pw string) (string, error) {
	pwHash, err := argon2id.CreateHash(pw, argon2id.DefaultParams)
	if err != nil {
		return "", fmt.Errorf("Error hashing password: %v", err)
	}
	return pwHash, nil
}

func CheckPasswordHash(pw, hash string) (bool, error) {
	same, err := argon2id.ComparePasswordAndHash(pw, hash)
	if err != nil {
		return false, fmt.Errorf("Error comparing password: %v", err)
	}

	return same, nil
}
