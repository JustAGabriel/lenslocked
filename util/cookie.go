package util

import "net/http"

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
