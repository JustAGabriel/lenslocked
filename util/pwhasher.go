package util

import "golang.org/x/crypto/bcrypt"

func GetPasswordHash(pw string) (hash string, err error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashBytes), nil
}
