package util

import "log"

const (
	SessionTokenLength = 32
)

func GetSessionToken() string {
	s, err := GetRandomBase64String(SessionTokenLength)
	if err != nil {
		log.Fatal(err)
	}

	return s
}
