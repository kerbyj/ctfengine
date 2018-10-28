package userapi

import (
	"ctfEngine/backend/common"
	"ctfEngine/backend/database"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"strconv"
)

func UserInfo(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*common.JwtCustomClaims)
	id := claims.UserId

	var request, error = database.DB.Query("SELECT COUNT(*), username, points, ROUND(flagfalse/(flagfalse+flagright), 2) as bff, commandid, FIND_IN_SET( points, ("+
		"SELECT GROUP_CONCAT( points "+
		"ORDER BY points DESC ) "+
		"FROM users ) "+
		") AS rank "+
		"FROM users where id = ?", id)

	if error != nil {
		log.Println(error)
	}
	defer request.Close()

	var countCheck, points, rank, command, captainId int
	var username, commandName string
	var bffactor float64

	request.Next()
	request.Scan(&countCheck, &username, &points, &bffactor, &command, &rank)

	log.Println("command", command)
	if command == 0 {
		commandName = "no command :("
	} else {
		var requestCommandName, errGetCommand = database.DB.Query("SELECT name, captainid FROM command WHERE id=?", command)
		if errGetCommand != nil {
			log.Println(errGetCommand)
		}

		requestCommandName.Next()
		requestCommandName.Scan(&commandName, &captainId)
		if captainId == id {
			commandName = "Captain of " + commandName
		}
	}

	if countCheck == 0 {
		return c.String(http.StatusBadRequest, "not found")
	}

	var dataOut = map[string]string{
		"name":          username,
		"command":       commandName,
		"Points":        strconv.Itoa(points),
		"Bruteforcer factor": fmt.Sprintf("%.2f%%", bffactor),
		"Overall place": strconv.Itoa(rank),
	}

	return c.JSON(http.StatusOK, dataOut)
}

func UserInfoByParameter(c echo.Context) error {
	var requestedUser = c.Param("name")

	var request = database.DB.QueryRow("SELECT COUNT(*), id, email, status FROM users WHERE username=?", requestedUser)

	var countCheck, id, status int
	var email, command string

	request.Scan(&countCheck, &id, &email, &status)

	if countCheck == 0 {
		return c.String(http.StatusBadRequest, "not found")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"id":      strconv.Itoa(id),
		"command": command,
		"status":  "ok",
	})
}

type TopUserOut struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Points   int    `json:"points"`
	Command  string `json:"command"`
}

func TopUserForAlltime(c echo.Context) error {
	var request, errorGetTop = database.DB.Query("SELECT users.id, username, points, command.name FROM users left join command on users.commandid = command.id ORDER BY points DESC LIMIT 50")

	if errorGetTop != nil {
		log.Print(errorGetTop)
		return c.String(http.StatusConflict, "take my bear, i need to fix something")
	}
	defer request.Close()

	var id, points int
	var username, command string
	var outData []TopUserOut

	for request.Next() {
		request.Scan(&id, &username, &points, &command)

		if command == "" {
			command = "Without command"
		}

		outData = append(outData, TopUserOut{
			id,
			username,
			points,
			command,
		})
	}

	return c.JSON(http.StatusOK, outData)
}
