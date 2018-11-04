package admin

import (
	"ctfEngine/backend/common"
	"ctfEngine/backend/database"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"net/http"
)

func CheckRights(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*common.JwtCustomClaims)
	id := claims.UserId

	var checkStatus, errCheckStatus = database.DB.Query("select status from users where id=?", id)

	if errCheckStatus != nil {
		return c.String(http.StatusOK, "we have a problem")
	}

	var admin int
	checkStatus.Next()
	checkStatus.Scan(&admin)

	if admin == 0 {
		return c.HTML(http.StatusOK, "<script>location.replace('/')</script>")
	} else {
		return c.String(http.StatusOK, "welcome admin")
	}

}