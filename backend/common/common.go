package common

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

// jwtCustomClaims are custom claims extending default ones.
type JwtCustomClaims struct { // TODO Вынести в common нормально
	Name   string `json:"name"`
	Admin  int    `json:"admin"`
	UserId int    `json:"id"`
	jwt.StandardClaims
}

func CreateCookie(name string, value string, httpOnly bool, path string, expires time.Time) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = value
	cookie.Path = path
	cookie.Expires = expires
	cookie.HttpOnly = httpOnly

	return cookie
}

const (
	JWT_AUTH_KEY = "user"
)
