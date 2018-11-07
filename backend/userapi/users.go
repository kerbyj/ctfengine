package userapi

import (
	"crypto/md5"
	"ctfEngine/backend/common"
	"ctfEngine/backend/database"
	"encoding/hex"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"html"
	"log"
	"net/http"
	"strconv"
)

func UserInfo(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*common.JwtCustomClaims)
	id := claims.UserId

	var request, error = database.DB.Query("SELECT COUNT(*), username, IFNULL((SELECT sum(points) FROM ctfengine.rating where userid=2 and team = 0), 0), IFNULL(ROUND(flagfalse/(flagfalse+flagright), 2), 0) as bff, commandid FROM users where id = ?", id)

	if error != nil {
		log.Println(error)
	}
	defer request.Close()

	var countCheck, points, command, captainId int
	var username, commandName string
	var bffactor float64

	request.Next()
	request.Scan(&countCheck, &username, &points, &bffactor, &command)
	//log.Println(username, points, bffactor, command)

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
		"name":               html.EscapeString(username),
		"command":            html.EscapeString(commandName),
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
			Name:          html.EscapeString(username),
			Points:        points,
			Solved:        solved,
			AllTasksCount: taskcount,
			Place:         counterPlace,
		})
	}
	return c.JSON(http.StatusOK, dataOut)
	//return c.String(http.StatusOK, "Ooops. We have a problem")
}

func ChangeUsername(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*common.JwtCustomClaims)
	id := claims.UserId

	var (
		newName = c.FormValue("newName")

		countName int
	)

	var getCountForUsername, errorGetCountForUsername = database.DB.Query("SELECT count(*) FROM users WHERE username=?", newName)
	if errorGetCountForUsername != nil {
		log.Println(errorGetCountForUsername, "for", id)
		return c.String(http.StatusServiceUnavailable, "Ooops")
	}
	getCountForUsername.Next()
	getCountForUsername.Scan(&countName)

	if countName == 1 {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "error",
			"error":  "the name is already taken",
		})
	}

	database.DB.Query("update users set username=? where id=?", newName, id)
	return c.JSON(http.StatusOK, map[string]string{
		"status": "success",
	})
}

func ChangePassword(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*common.JwtCustomClaims)
	id := claims.UserId

	var (
		oldPassword = c.FormValue("oldPassword")
		newPassword = c.FormValue("newPassword")

		passwordFromDb string
	)

	var getPasswordForUser, errorGetPassword = database.DB.Query("SELECT password FROM users WHERE id=?", id)
	if errorGetPassword != nil {
		log.Println(errorGetPassword, "for", id)
		return c.String(http.StatusServiceUnavailable, "Ooops")
	}
	getPasswordForUser.Next()
	getPasswordForUser.Scan(&passwordFromDb)

	var hashOldPass = md5.New()
	hashOldPass.Write([]byte(oldPassword))
	var oldPassHashString = hex.EncodeToString(hashOldPass.Sum(nil))

	if passwordFromDb == oldPassHashString {
		var hashForNewPassword = md5.New()
		//io.WriteString(hashForNewPassword, newPassword)
		hashForNewPassword.Write([]byte(newPassword))
		var hashForNewPasswordString = hex.EncodeToString(hashForNewPassword.Sum(nil))

		var _, errorUpdatePassword = database.DB.Query("UPDATE users SET password=? WHERE id=?", hashForNewPasswordString, id)
		if errorUpdatePassword != nil {
			return c.JSON(http.StatusOK, map[string]string{
				"status": "error",
				"error":  "some problems with db",
			})
		}

		return c.JSON(http.StatusOK, map[string]string{
			"status": "success",
		})

	} else {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "error",
			"error":  "invalid password",
		})
	}
}

func GetCommandInfoForSettings(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*common.JwtCustomClaims)
	id := claims.UserId

	var requestCommandId, errorGetCommandId = database.DB.Query("select commandid, command.captainid, command.name, command.invite from users left join command on users.commandid = command.id where users.id=?", id)
	if errorGetCommandId != nil {
		log.Println(errorGetCommandId, "on id", id)
		return c.JSON(http.StatusOK, map[string]string{
			"status": "error",
			"error":  "some problems with db",
		})
	}

	var (
		commandId, captainId       int
		commandName, commandInvite string
	)
	requestCommandId.Next()
	requestCommandId.Scan(&commandId, &captainId, &commandName, &commandInvite)

	if commandId == 0 {
		return c.JSON(http.StatusOK, map[string]int{
			"command_id": 0,
		})
	} else {
		var getCommandList, errorGetCommandList = database.DB.Query("select id, username from users where commandid=?", commandId)
		if errorGetCommandList != nil {
			log.Println(errorGetCommandList)
			return c.JSON(http.StatusOK, map[string]string{
				"status": "error",
				"error":  "some problems with db",
			})
		}

		var (
			commandMemberId       int
			commandMemberUsername string
			memberList            = make(map[int]string)
		)

		type commandDataOut struct {
			YourId        int            `json:"your_id"`
			CaptainId     int            `json:"captain_id"`
			CommandId     int            `json:"command_id"`
			CommandName   string         `json:"command_name"`
			CommandInvite string         `json:"command_invite"`
			Members       map[int]string `json:"members"`
		}

		for getCommandList.Next() {
			getCommandList.Scan(&commandMemberId, &commandMemberUsername)
			memberList[commandMemberId] = commandMemberUsername
		}

		return c.JSON(http.StatusOK, commandDataOut{
			YourId:        id,
			CaptainId:     captainId,
			CommandId:     commandId,
			CommandName:   commandName,
			Members:       memberList,
			CommandInvite: commandInvite,
		})
	}
}

func LeaveCommand(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*common.JwtCustomClaims)
	id := claims.UserId

	var _, errorDeleteCommand = database.DB.Query("UPDATE users SET commandid = 0 WHERE id = ?", id)
	if errorDeleteCommand != nil {
		log.Println(errorDeleteCommand, "for", id)
		return c.JSON(http.StatusOK, map[string]string{
			"status": "error",
		})
	}
	return c.JSON(http.StatusOK, map[string]string{
		"status": "success",
	})
}

func CreateCommand(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*common.JwtCustomClaims)
	id := claims.UserId

	var commandName = c.FormValue("commandname")

	var checkNameExist, errorCheckNameExist = database.DB.Query("select (select count(*) from command where name = ?), (select commandid from users where id=?)", commandName, id)
	if errorCheckNameExist != nil {
		log.Println(errorCheckNameExist, "for", commandName)
		return c.JSON(http.StatusOK, map[string]string{
			"status": "error",
		})
	}

	var commandStatusExist, alreadyInCommand int
	checkNameExist.Next()
	checkNameExist.Scan(&commandStatusExist, &alreadyInCommand)

	if commandStatusExist != 0 {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "not available",
		})
	}

	if alreadyInCommand != 0 {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "already in command",
		})
	}

	var _, errorCreateCommand = database.DB.Query("INSERT into command (name, captainid) VALUES (?, ?)", commandName, id)
	if errorCreateCommand != nil {
		log.Println(errorCreateCommand, "for", commandName, id)
		return c.JSON(http.StatusOK, map[string]string{
			"status": "error",
		})
	}

	var (
		commandId int
	)

	var commandIdRequest, _ = database.DB.Query("select id from command where captainid=?", id)
	commandIdRequest.Next()
	commandIdRequest.Scan(&commandId)

	var _, errorUpdateCommandStatusForCaptain = database.DB.Query("update users set commandid=? where id=?", commandId, id)
	if errorUpdateCommandStatusForCaptain != nil {
		log.Println(errorUpdateCommandStatusForCaptain, "for", commandId, id)
		return c.JSON(http.StatusOK, map[string]string{
			"status": "error",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": "success",
	})
}

func RenameCommand(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*common.JwtCustomClaims)
	id := claims.UserId

	var commandForRename = c.FormValue("commandname")

	var getCommandId, errorGetCaptainId = database.DB.Query("select (select id from command where captainid=?),(select count(*) from command where name=?)", id, commandForRename)
	if errorGetCaptainId != nil {
		log.Println(errorGetCaptainId, "for", commandForRename, id)
		return c.JSON(http.StatusOK, map[string]string{
			"status": "error",
		})
	}
	var commandId, checkName int
	getCommandId.Next()
	getCommandId.Scan(&commandId, &checkName)

	if checkName != 0 {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "Name is taken",
		})
	}

	if commandId != 0 {
		database.DB.Query("update command set name=? where id=?", commandForRename, commandId)
		return c.JSON(http.StatusOK, map[string]string{
			"status": "success",
		})
	}
	return c.JSON(http.StatusOK, map[string]string{
		"status": "you are not captain",
	})
}

func DeleteCommand(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*common.JwtCustomClaims)
	id := claims.UserId

	var requestedCommandForDelete = c.FormValue("commandid")

	var checkCaptainId, errorCheckCaptainId = database.DB.Query("select id, captainid from command where id=?", requestedCommandForDelete)
	if errorCheckCaptainId != nil {
		log.Println(errorCheckCaptainId)
		return c.JSON(http.StatusOK, map[string]string{
			"status": "error",
		})
	}

	var commandId, captainId int
	checkCaptainId.Next()
	checkCaptainId.Scan(&commandId, &captainId)

	//log.Println(captainId, id)
	if captainId != id {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "error",
			"error":  "you are not captain",
		})
	}

	database.DB.Query("update users set commandid=0 where commandid=?", commandId)
	database.DB.Query("delete from command where id=?", commandId)

	return c.JSON(http.StatusOK, map[string]string{
		"status": "success",
	})
}

func DropUserFromCommand(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*common.JwtCustomClaims)
	id := claims.UserId

	var requestedUserForDelete = c.FormValue("userid")

	var checkCaptainId, errorCheckCaptainId = database.DB.Query("select id, count(*) from command where captainid=?", id)
	if errorCheckCaptainId != nil {
		log.Println(errorCheckCaptainId)
		return c.JSON(http.StatusOK, map[string]string{
			"status": "error",
		})
	}

	var commandId, checkCount, currentUserCommand int

	checkCaptainId.Next()
	checkCaptainId.Scan(&commandId, &checkCount)

	if checkCount == 0 {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "error",
			"error":  "you are not captain",
		})
	}

	var checkCurrentCommand, errorCheckCurrentCommand = database.DB.Query("select commandid from users where id=?", requestedUserForDelete)
	if errorCheckCurrentCommand != nil {
		log.Println(errorCheckCurrentCommand)
		return c.JSON(http.StatusOK, map[string]string{
			"status": "error",
		})
	}

	checkCurrentCommand.Next()
	checkCurrentCommand.Scan(&currentUserCommand)

	if currentUserCommand != commandId {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "error",
			"error":  "kysh otsyda",
		})
	}

	database.DB.Query("update users set commandid=0 where id=?", requestedUserForDelete)

	return c.JSON(http.StatusOK, map[string]string{
		"status": "success",
	})
}

func JoinCommandViaInvite(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*common.JwtCustomClaims)
	id := claims.UserId

	var inviteCode = c.FormValue("invite")
	var getCommandCreds, errorGetCommandCreds = database.DB.Query("select id, count(*) from command where invite=?", inviteCode)
	if errorGetCommandCreds != nil {
		log.Println(errorGetCommandCreds)
		return c.JSON(http.StatusOK, map[string]string{
			"status": "error",
		})
	}

	var commandId, countCommandCheck int
	getCommandCreds.Next()
	getCommandCreds.Scan(&commandId, &countCommandCheck)

	if countCommandCheck != 1 {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "error",
			"error":  "bad invite link",
		})
	}

	database.DB.Query("UPDATE users SET commandid=? WHERE id=?", commandId, id)
	return c.JSON(http.StatusOK, map[string]string{
		"status": "success",
	})
}
