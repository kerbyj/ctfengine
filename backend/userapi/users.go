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

	var request, error = database.DB.Query("SELECT COUNT(*), username, (SELECT sum(points) FROM ctfengine.rating where userid=2 and team = 0), ROUND(flagfalse/(flagfalse+flagright), 2) as bff, commandid FROM users where id = ?", id)

	if error != nil {
		log.Println(error)
	}
	defer request.Close()

	var countCheck, points, command, captainId int
	var username, commandName string
	var bffactor float64

	request.Next()
	request.Scan(&countCheck, &username, &points, &bffactor, &command)

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
		"name":               username,
		"command":            commandName,
		"Points":             strconv.Itoa(points),
		"Bruteforcer factor": fmt.Sprintf("%.2f%%", bffactor),
		//"Overall place": strconv.Itoa(rank),
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

func GetTopForContest(c echo.Context) error {
	var requestedContest = c.Param("contestid")
	var checkContestExist, errorCheckContest = database.DB.Query("SELECT COUNT(*), name, type FROM contests WHERE id=?", requestedContest)
	if errorCheckContest != nil {
		log.Println(errorCheckContest)
		return c.String(http.StatusOK, "Ooops. We have a problem in checkContestExist")
	}
	var (
		counter                  int
		contestName, contestType string
	)
	checkContestExist.Next()
	checkContestExist.Scan(&counter, &contestName, &contestType)

	if counter == 0 {
		return c.String(http.StatusOK, "not found")
	}

	type topForContest struct {
		Name          string `json:"name"`
		Points        int    `json:"points"`
		Solved        int    `json:"solved"`
		AllTasksCount int    `json:"all_tasks_count"`
		Place         int    `json:"place"`
	}

	var requestQuery string

	if contestType == "alone" {
		requestQuery = "SELECT userid, rating.points, " +
			"(select count(*) from pwnedby where userid=rating.userid and contestid=?) as solved, " +
			"(select count(*) from tasks where contestid=?) as taskscount, " +
			"username from rating " +
			"left join users ON rating.userid=users.id where contest_id=? order by points desc"

	} else if contestType == "team" {
		requestQuery = "SELECT team, rating.points, " +
			"(select count(*) from pwnedby where command_id=rating.team and contestid=?) as solved, " +
			"(select count(*) from tasks where contestid=?) as taskscount, " +
			"name from rating " +
			"left join command ON rating.team=command.id where contest_id=? order by points desc"

	}
	var getTopForContest, errGetTopForContest = database.DB.Query(requestQuery, requestedContest, requestedContest, requestedContest)

	if errGetTopForContest != nil {
		log.Println(errorCheckContest)
		return c.String(http.StatusOK, "Ooops. We have a problem")
	}
	var (
		userid, points, solved, taskcount int
		username                          string

		counterPlace = 0

		dataOut []topForContest
	)

	for getTopForContest.Next() {
		getTopForContest.Scan(&userid, &points, &solved, &taskcount, &username)
		counterPlace++

		dataOut = append(dataOut, topForContest{
			Name:          username,
			Points:        points,
			Solved:        solved,
			AllTasksCount: taskcount,
			Place:         counterPlace,
		})
	}
	return c.JSON(http.StatusOK, dataOut)
	//return c.String(http.StatusOK, "Ooops. We have a problem")
}
