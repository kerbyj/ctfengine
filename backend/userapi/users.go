package userapi

import (
	"ctfEngine/backend/common"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)



func UserInfo(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*common.JwtCustomClaims)
	name := claims.Name
	id := claims.UserId

	var dataOut = map[string]string{
		"name": name,
		"id": strconv.Itoa(id),
		"command": "test",
	}

	return c.JSON(http.StatusOK, dataOut)
}
