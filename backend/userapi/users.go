package userapi

import (
	"ctfEngine/backend/common"
	"ctfEngine/backend/database"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"strconv"
)

func UserInfo(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*common.JwtCustomClaims)
	name := claims.Name
	id := claims.UserId

	log.Println(name, id)

	var request, error = database.DB.Query("SELECT COUNT(*), username, points, command, FIND_IN_SET( points, ("+
		"SELECT GROUP_CONCAT( points "+
		"ORDER BY points DESC ) "+
		"FROM users ) "+
		") AS rank "+
		"FROM users where id = ?", id)

	if error != nil {
		log.Println(error)
	}
	defer request.Close()

	var countCheck, points, rank int
	var username, command string

	request.Next()
	request.Scan(&countCheck, &username, &points, &command, &rank)

	if countCheck == 0 {
		return c.String(http.StatusBadRequest, "not found")
	}

	if command == "" {
		command = "¯\\_(ツ)_/¯"
	}

	var dataOut = map[string]string{
		"name":          username,
		"command":       command,
		"Points":        strconv.Itoa(points),
		"Overall place": strconv.Itoa(rank),
	}

	return c.JSON(http.StatusOK, dataOut)
}

func UserInfoByParameter(c echo.Context) error {
	var requestedUser = c.Param("name")

	var request = database.DB.QueryRow("SELECT COUNT(*), id, email, command, status FROM users WHERE username=?", requestedUser)

	var countCheck, id, status int
	var email, command string

	request.Scan(&countCheck, &id, &email, &command, &status)

	if countCheck == 0 {
		return c.String(http.StatusBadRequest, "not found")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"id":      strconv.Itoa(id),
		"command": command,
		"status":  "ok",
	})
}
