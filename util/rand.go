package util

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func getRandomBytes(length int) ([]byte, error) {
	bytesArr := make([]byte, length)
	generatedBytes, err := rand.Read(bytesArr)
	if err != nil {
		return nil, err
	}

	if generatedBytes != length {
		err := fmt.Errorf("could not generate expected amount of random bytes (generated: %d, expected: %d)", generatedBytes, length)
		return nil, err
	}

	return bytesArr, nil
}

func GetRandomBase64String(length int) (string, error) {
	randomBytes, err := getRandomBytes(length)
	if err != nil {
		return "", err
	}

	base64String := base64.URLEncoding.EncodeToString(randomBytes)
	return base64String, nil
}
