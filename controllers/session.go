package controllers

import (
	"log"

	"github.com/justagabriel/lenslocked/util"
)

const (
	SessionTokenLength = 32
)

func GetSessionToken() string {
	s, err := util.GetRandomBase64String(SessionTokenLength)
	if err != nil {
		log.Fatal(err)
	}

	return s
}
