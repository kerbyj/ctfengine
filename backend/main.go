package main

import (
	"database/sql"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

var (
	database *sql.DB
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func loadDB() {
	db, err := sql.Open("sqlite3", "db2.sqlite")
	checkErr(err)
	database = db
	log.Println("Database start")
}

func accessible(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}

func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	name := claims.Name
	return c.String(http.StatusOK, "Welcome "+name+"!")
}

func main() {
	loadDB() // Load db

	e := echo.New()

	// Middleware
	//e.Use(middleware.Logger())
	//e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"}, // TODO Debug mode! Change on real domain
		AllowMethods: []string{echo.GET, echo.POST},
	}))



	// Login route
	e.POST("/login", login)

	// register root
	e.POST("/register", register)

	// Unauthenticated route
	e.GET("/", accessible)

	// Restricted group
	r := e.Group("/restricted")

	// Configure middleware with the custom claims type
	config := middleware.JWTConfig{
		Claims:     &jwtCustomClaims{},
		SigningKey: secretJWTkey,
	}
	r.Use(middleware.JWTWithConfig(config))
	r.GET("", restricted)

	e.Logger.Fatal(e.Start(":80"))
}
