package util

import (
	"fmt"
	"net/http"
)

func SetCookie(r http.ResponseWriter, name string, value string) error {
	c := http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: false,
	}

	if err := c.Valid(); err != nil {
		return err
	}

	http.SetCookie(r, &c)
	return nil
}

func GetSessionTokenFromCookie(cookieName string, request *http.Request) (string, error) {
	c, err := request.Cookie(cookieName)
	if err != nil {
		return "", fmt.Errorf("error while trying to read cookie: %+v", err)
	}

	return c.Value, nil
}
