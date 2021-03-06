package main

import (
	"github.com/kerbyj/ctfengine/backend/admin"
	"github.com/kerbyj/ctfengine/backend/boardapi"
	"github.com/kerbyj/ctfengine/backend/common"
	"github.com/kerbyj/ctfengine/backend/database"
	"github.com/kerbyj/ctfengine/backend/taskapi"
	"github.com/kerbyj/ctfengine/backend/userapi"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
	"time"

	//_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
)

var (
	executionPath string
)

func customHTTPErrorHandler(err error, c echo.Context) {
	/*
		code := http.StatusInternalServerError
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}

		if code == 400 {
			c.HTML(http.StatusOK, "<script>location.replace('/login')</script>")
		}
	*/
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

	var files = e.Group("/files/")
	fs := http.FileServer(http.Dir(executionPath + "/files/"))
	files.GET("*", echo.WrapHandler(http.StripPrefix("/files", fs)))

	// Middleware
	//e.Use(middleware.Logger())
	//e.Use(middleware.Secure())
	//e.Pre(middleware.HTTPSRedirect())

	e.HTTPErrorHandler = customHTTPErrorHandler

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"}, // TODO Debug mode! Change on real domain
		AllowMethods: []string{echo.GET, echo.POST},
	}))

	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLength:  32,
		TokenLookup:  "header:" + echo.HeaderXCSRFToken,
		ContextKey:   "csrf",
		CookieName:   "_csrf",
		CookieMaxAge: 86400,
	}))

	e.POST("/api/auth/login", login)       // Login user & Create JWT
	e.POST("/api/auth/register", register) // Register user & Create JWT
	e.GET("/api/user/:name", userapi.UserInfoByParameter)

	e.GET("/", func(c echo.Context) error {
		return c.File(executionPath + "/frontend/hello.html")
	}) // Default page with landing

	e.GET("/login", func(c echo.Context) error {
		return c.File(executionPath + "/frontend/login.html")
	}) // Login/Register page

	e.GET("/logout", func(c echo.Context) error {
		c.SetCookie(common.CreateCookie("token", "", true, "/", time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC)))

		return c.HTML(http.StatusOK, "<script>location.replace('/login')</script>") // TODO not work properly
	})

	config := middleware.JWTConfig{
		Claims:      &common.JwtCustomClaims{},
		SigningKey:  secretJWTkey,
		TokenLookup: "cookie:token",
	}

	b := e.Group("/board")
	t := e.Group("/tasks")
	s := e.Group("/scoreboard")
	c := e.Group("/settings")
	adm := e.Group("/admin")
	api := e.Group("/api")

	b.Use(middleware.JWTWithConfig(config))
	t.Use(middleware.JWTWithConfig(config))
	s.Use(middleware.JWTWithConfig(config))
	c.Use(middleware.JWTWithConfig(config))
	adm.Use(middleware.JWTWithConfig(config))
	api.Use(middleware.JWTWithConfig(config))

	b.GET("", func(c echo.Context) error { // b /board
		return c.File(executionPath + "/frontend/board.html")
	})

	t.GET("", func(c echo.Context) error { // t /tasks
		return c.File(executionPath + "/frontend/tasks.html")
	})

	s.GET("", func(c echo.Context) error { // s scoreboard
		return c.File(executionPath + "/frontend/scoreboard.html")
	})

	c.GET("", func(c echo.Context) error { // c settings
		return c.File(executionPath + "/frontend/settings.html")
	})

	adm.GET("", admin.MainPage)
	adm.POST("/createContest", admin.CreateContest)
	adm.POST("/createTask", admin.CreateTask)

	api.GET("/user/info", userapi.UserInfo)                    // Get info for logged-in user
	api.GET("/users/topForAllTime", userapi.TopUserForAlltime) // For scoreboard
	api.GET("/users/getTopForContest/:contestid", userapi.GetTopForContest)
	api.POST("/users/ChangePassword", userapi.ChangePassword)
	api.POST("/users/ChangeUsername", userapi.ChangeUsername)

	api.GET("/users/getCommandStatusForSettings", userapi.GetCommandInfoForSettings)
	api.POST("/users/LeaveCommand", userapi.LeaveCommand)
	api.POST("/users/CreateCommand", userapi.CreateCommand)
	api.POST("/users/RenameCommand", userapi.RenameCommand)
	api.POST("/users/DeleteCommand", userapi.DeleteCommand)
	api.POST("/users/DropUserFromCommand", userapi.DropUserFromCommand)
	api.POST("/users/JoinCommandViaInvite", userapi.JoinCommandViaInvite)

	api.GET("/tasks/getAlwaysAliveTasks", taskapi.GetAlwaysAliveTasks) //
	api.GET("/tasks/getContestList", taskapi.GetContestList)
	api.GET("/tasks/getTaskById/:id", taskapi.GetTaskById)
	api.GET("/tasks/getFile/:id", taskapi.GetTaskById)
	api.POST("/tasks/checkFlag", taskapi.CheckFlag)
	//api.GET("/tasks/GetContestTasks", )

	api.GET("/board/getstats", boardapi.BoardStats)

	//e.Logger.Fatal(e.StartTLS(":443", "cert.pem", "key.pem"))
	e.Logger.Fatal(e.Start(":1080"))
}
