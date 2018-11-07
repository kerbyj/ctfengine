package main

import (
	"crypto/md5"
	"ctfEngine/backend/common"
	"ctfEngine/backend/database"
	"encoding/hex"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"time"
)

var (
	secretJWTkey = []byte("1234")
)

const (
	salt = ""
)

func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	rows := database.DB.QueryRow("SELECT COUNT(*), id, password, status FROM users WHERE username=?", username)

	var count, id int
	var passwordhash string
	var status int

	rows.Scan(&count, &id, &passwordhash, &status)

	if count == 0 {
		return echo.ErrUnauthorized
	}

	hasher := md5.New()
	hasher.Write([]byte(salt + password))

	if passwordhash == hex.EncodeToString(hasher.Sum(nil)) {
		// Set custom claims
		claims := &common.JwtCustomClaims{
			username,
			status,
			id,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 12).Unix(),
			},
		}

		// Create token with claims
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Generate encoded token and send it as response.
		t, err := token.SignedString(secretJWTkey)
		if err != nil {
			return c.JSON(http.StatusOK, map[string]string{
				"status": "error",
				"error":  "error inside core",
			})
		}

		c.SetCookie(common.CreateCookie("token", t, true, "/", time.Now().Add(12*time.Hour)))

		return c.JSON(http.StatusOK, map[string]string{
			"status": "success",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": "error",
		"error":  "invalid credentials",
	})
}

func register(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	email := c.FormValue("email")

	rows := database.DB.QueryRow("SELECT(SELECT COUNT(*) FROM users WHERE username=?),(SELECT COUNT(*) FROM users WHERE email=?)", username, email)

	var countByUsers, countByEmails int

	rows.Scan(&countByUsers, &countByEmails)

	if countByUsers+countByEmails != 0 {
		log.Println("Counts not null")
		return echo.ErrUnauthorized
	}
	hasher := md5.New()
	hasher.Write([]byte(salt + password))
	var passwordHash = hex.EncodeToString(hasher.Sum(nil))
	stmt, err := database.DB.Prepare("INSERT INTO users(username, password, email) values(?,?,?)")

	res, errExecInsert := stmt.Exec(username, passwordHash, email)
	stmt.Close()
	if errExecInsert != nil {
		log.Println(errExecInsert, "ErrExecInsert")
		return echo.ErrUnauthorized
	}
	var userId, _ = res.LastInsertId()
	// Set custom claims
	claims := &common.JwtCustomClaims{
		username,
		0,
		int(userId),
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 12).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString(secretJWTkey)
	if err != nil {
		return err
	}

	c.SetCookie(common.CreateCookie("token", t, true, "/", time.Now().Add(12*time.Hour)))

	return c.JSON(http.StatusOK, "ok")
}
