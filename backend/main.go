package main

import (
	"ctfEngine/backend/common"
	"ctfEngine/backend/database"
	"ctfEngine/backend/taskapi"
	"ctfEngine/backend/userapi"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
)

var (
	executionPath string
)

func landing(c echo.Context) error {
	return c.File(executionPath + "/frontend/hello.html")
}

func board(c echo.Context) error {
	return c.File(executionPath + "/frontend/board.html")
}

func tasks(c echo.Context) error {
	return c.File(executionPath + "/frontend/tasks.html")
}

func scoreboard(c echo.Context) error {
	return c.File(executionPath + "/frontend/scoreboard.html")
}

func settings(c echo.Context) error {
	return c.File(executionPath + "/frontend/settings.html")
}

func loginpage(c echo.Context) error {
	return c.File(executionPath + "/frontend/login.html")
}

func main() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	executionPath = filepath.Dir(ex)

	database.LoadDB(executionPath) // Load db

	e := echo.New()

	e.Static("/css", executionPath+"/frontend/css")
	e.Static("/js", executionPath+"/frontend/js")

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
	e.GET("/api/user/:name", userapi.UserInfoByParameter)

	// End API methods

	/*
		Page routing
	*/
	// Pages without checking user auth
	e.GET("/", landing)        // Default page with landing
	e.GET("/login", loginpage) // Login/Register page

	// Restricted group, need to get JWT

	// JWT config
	config := middleware.JWTConfig{
		Claims:      &common.JwtCustomClaims{},
		SigningKey:  secretJWTkey,
		TokenLookup: "cookie:token",
	}

	b := e.Group("/board")
	b.Use(middleware.JWTWithConfig(config))
	b.GET("", board) // Dashboard with stats

	t := e.Group("/tasks")
	t.Use(middleware.JWTWithConfig(config))
	t.GET("", tasks) // Tasks

	s := e.Group("/scoreboard")
	s.Use(middleware.JWTWithConfig(config))
	s.GET("", scoreboard) // Tasks

	с:= e.Group("/settings")
	с.Use(middleware.JWTWithConfig(config))
	с.GET("", settings) // Tasks

	api := e.Group("/api")
	api.Use(middleware.JWTWithConfig(config))
	api.GET("/users/info", userapi.UserInfo) // Get info for logged-in user
	api.GET("/users/topForAllTime", userapi.TopUserForAlltime) // For scoreboard

	api.GET("/tasks/getAlwaysAliveTasks", taskapi.GetAlwaysAliveTasks) //

	//api.POST("/changecommand", )
	//r.GET("/top", )

	e.Logger.Fatal(e.Start(":80"))
}
