package admin

import (
	"ctfEngine/backend/common"
	"ctfEngine/backend/database"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"net/http"
)

func CheckRights(id int) (int, error) {

	var checkStatus, errCheckStatus = database.DB.Query("select status from users where id=?", id)

	if errCheckStatus != nil {
		//return c.String(http.StatusOK, "we have a problem")
		return 0, errors.New("db problem")
	}

	var admin int
	checkStatus.Next()
	checkStatus.Scan(&admin)

	return admin, nil
}

func MainPage(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*common.JwtCustomClaims)
	id := claims.UserId

	var admin, errorCheckAdmin = CheckRights(id)

	if errorCheckAdmin != nil {
		return c.String(http.StatusOK, "we have a problem")
	}

	if admin == 0 {
		return c.HTML(http.StatusOK, "<script>location.replace('/')</script>")
	} else {
		//return c.String(http.StatusOK, "welcome admin")
		return c.File("frontend/admin.html")
	}
}

func CreateContest(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*common.JwtCustomClaims)
	id := claims.UserId

	var admin, errorCheckAdmin = CheckRights(id)

	if errorCheckAdmin != nil {
		return c.String(http.StatusOK, "we have a problem")
	}

	if admin == 0 {
		return c.String(http.StatusOK, "not admin")
	}

	var (
		contestName       = c.FormValue("contest_name")
		contestType       = c.FormValue("contest_type")
		contestVisibility = c.FormValue("visibility")
		contestPermit     = c.FormValue("permit")
	)

	var cc, _ = database.DB.Query("insert into contests (name, type, visibility, permit) values (?,?,?,?)", contestName, contestType, contestVisibility, contestPermit)
	defer cc.Close()

	return c.JSON(http.StatusOK, map[string]string{
		"status": "success",
	})
}

func CreateTask(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*common.JwtCustomClaims)
	id := claims.UserId

	var admin, errorCheckAdmin = CheckRights(id)

	if errorCheckAdmin != nil {
		return c.String(http.StatusOK, "we have a problem")
	}
	if admin == 0 {
		return c.String(http.StatusOK, "not admin")
	}

	var (
		taskName        = c.FormValue("task_name")
		taskFlag        = c.FormValue("task_flag")
		taskPrice       = c.FormValue("task_price")
		taskCategory    = c.FormValue("task_category")
		taskDescription = c.FormValue("task_description")
		taskContestId   = c.FormValue("task_contest")
	)

	var cc, _ = database.DB.Query("insert into tasks (name, contestid, value, flag, description, category) values (?, ?, ?, ?, ?, ?)", taskName, taskContestId, taskPrice, taskFlag, taskDescription, taskCategory)
	defer cc.Close()

	return c.JSON(http.StatusOK, map[string]string{
		"status": "success",
	})
}
