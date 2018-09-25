package main

import (
	"database/sql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/mattn/go-sqlite3"
	"log"
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

func landing(c echo.Context) error {
	return c.File("C:/users/kerby/go/src/ctfengine/frontend/hello.html")
}

func board(c echo.Context) error {
	//user := c.Get("user").(*jwt.Token)
	//claims := user.Claims.(*jwtCustomClaims)
	//name := claims.Name

	//return c.String(http.StatusOK, "Welcome "+name+"!")
	return c.File("C:/users/kerby/go/src/ctfengine/frontend/board.html")
}

func loginpage(c echo.Context) error {
	return c.File("C:/users/kerby/go/src/ctfengine/frontend/login.html")
}

func main() {
	loadDB() // Load db

	e := echo.New()

	e.Static("/css", "C:/users/kerby/go/src/ctfengine/frontend/css")

	// Middleware
	//e.Use(middleware.Logger())
	//e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"}, // TODO Debug mode! Change on real domain
		AllowMethods: []string{echo.GET, echo.POST},
	}))

	/*
	API methods
	*/

	e.POST("/api/auth/login", login)       // Login user & Create JWT
	e.POST("/api/auth/register", register) // Register user & Create JWT

	/*
	End API methods
	*/

	/*
	Page routing
	*/
	// Pages without checking user auth

	e.GET("/", landing) // Default page with landing
	e.GET("/login", loginpage) // Login/Register page

	// Restricted group, need to get JWT

	// JWT config
	config := middleware.JWTConfig{
		Claims:      &jwtCustomClaims{},
		SigningKey:  secretJWTkey,
		TokenLookup: "cookie:token",
	}

	r := e.Group("/board")
	r.Use(middleware.JWTWithConfig(config))
	r.GET("", board) // Dashboard with stats
	//r.GET("/top", )

	e.Logger.Fatal(e.Start(":80"))
}
