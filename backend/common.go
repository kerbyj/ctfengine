package main

import (
	"net/http"
	"time"
)

func createCookie(name string, value string, httpOnly bool, path string) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = value
	cookie.Path = path
	cookie.Expires = time.Now().Add(12 * time.Hour)
	cookie.HttpOnly = httpOnly

	return cookie
}
