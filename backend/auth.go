package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
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

// jwtCustomClaims are custom claims extending default ones.
type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin int    `json:"admin"`
	jwt.StandardClaims
}


func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	rows, err := database.Query("SELECT COUNT(*), password, status FROM users WHERE username==?", username)

	if err != nil {
		log.Println(err)
	}

	var count int
	var passwordhash string
	var status int

	rows.Next()
	rows.Scan(&count, &passwordhash, &status)
	rows.Close()

	if count == 0 {
		return echo.ErrUnauthorized
	}

	hasher := md5.New()
	hasher.Write([]byte(salt + password))
	fmt.Println(hex.EncodeToString(hasher.Sum(nil)))

	if passwordhash == hex.EncodeToString(hasher.Sum(nil)) {
		// Set custom claims
		claims := &jwtCustomClaims{
			username,
			status,
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
		return c.JSON(http.StatusOK, echo.Map{
			"token": t,
		})
	}

	return echo.ErrUnauthorized
}

func register(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	email := c.FormValue("email")

	rows, err := database.Query("SELECT(SELECT COUNT(*) FROM users WHERE username==?),(SELECT COUNT(*) FROM users WHERE email==?)", username, email)

	if err != nil {
		log.Println(err, "ErrCountQuery")
	}

	var countByUsers, countByEmails int
	rows.Next()
	rows.Scan(&countByUsers, &countByEmails)
	rows.Close()

	if countByUsers + countByEmails != 0 {
		log.Println("Counts not null")
		return echo.ErrUnauthorized
	}

	hasher := md5.New()
	hasher.Write([]byte(salt + password))
	var passwordHash = hex.EncodeToString(hasher.Sum(nil))

	stmt, err := database.Prepare("INSERT INTO users(username, password, email) values(?,?,?)")
	checkErr(err)

	_, errExecInsert := stmt.Exec(username, passwordHash, email)
	stmt.Close()

	if errExecInsert != nil {
		log.Println(errExecInsert, "ErrExecInsert")
		return echo.ErrUnauthorized
	}

	// Set custom claims
	claims := &jwtCustomClaims{
		username,
		0,
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
	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}