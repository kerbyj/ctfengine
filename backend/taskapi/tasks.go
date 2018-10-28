package taskapi

import (
	"ctfEngine/backend/common"
	"ctfEngine/backend/database"
	"database/sql"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"runtime"
	"strconv"
)

type contestStruct struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	StartAt  string `json:"start_at"`
	FinishAt string `json:"finish_at"`
}

func GetContestList(c echo.Context) error {
	var getAllContests, errGetAllContests = database.DB.Query("SELECT * from contests")
	if errGetAllContests != nil {
		return c.String(http.StatusBadRequest, "shit happens")
	}

	var (
		dataOut           []contestStruct
		id                int
		name, start, stop string
	)

	for getAllContests.Next() {
		getAllContests.Scan(&id, &name, &start, &stop)
		dataOut = append(dataOut, contestStruct{
			Id:       id,
			Name:     name,
			StartAt:  start,
			FinishAt: stop,
		})
	}

	return c.JSON(http.StatusOK, dataOut)
}

type taskStruct struct {
	Name       string `json:"Name"`
	Value      int    `json:"value"`
	Id         int    `json:"id"`
	Solved     bool   `json:"solved"`
	SubmitTime string `json:"submit_time"`
	Contest    string `json:"contest"`
	Category   string `json:"category"`
}

func GetAlwaysAliveTasks(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*common.JwtCustomClaims)
	userid := claims.UserId

	var userCommandId int
	var getCommandId, _ = database.DB.Query("SELECT commandid FROM users WHERE id=?", userid)
	getCommandId.Next()
	getCommandId.Scan(&userCommandId)

	var GetAllTasks, errGetAliveTasks = database.DB.Query("SELECT tasks.id, category, tasks.name, value, contests.name FROM tasks left join contests on tasks.contestid = contests.id order by contestid")
	if errGetAliveTasks != nil {
		return c.String(http.StatusBadRequest, "shit happens")
	}
	var (
		alreadySolvedByUser = make(map[int]string)
		taskid              int
		time                string

		taskOut                     = make(map[string][]taskStruct)
		id, value                   int
		category, name, contestName string
	)

	var getAlreadySolvedTasks, _ = database.DB.Query("SELECT taskid, time FROM pwnedby WHERE userid=? OR command_id=?", userid, userCommandId)
	for getAlreadySolvedTasks.Next() {
		getAlreadySolvedTasks.Scan(&taskid, &time)
		alreadySolvedByUser[taskid] = time
	}

	var statusSolved bool
	var submitTime string

	for GetAllTasks.Next() {
		GetAllTasks.Scan(&id, &category, &name, &value, &contestName)

		if time, ok := alreadySolvedByUser[id]; ok {
			statusSolved = true
			submitTime = time
		} else {
			statusSolved = false
			submitTime = "no"
		}

		taskOut[contestName] = append(taskOut[contestName], taskStruct{
			Name:       name,
			Value:      value,
			Id:         id,
			Solved:     statusSolved,
			SubmitTime: submitTime,
			//Contest:    contestName,
			Category: category,
		})
		// Please, create issue if you have better way for this method
	}

	runtime.GC()
	return c.JSON(http.StatusOK, taskOut)
}

func GetTaskById(c echo.Context) error {
	var requestedTask = c.Param("id")
	var request, errGetTask = database.DB.Query("SELECT id, name, value, description FROM tasks WHERE id=?", requestedTask)

	if errGetTask != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	var id, value int
	var name, description string

	request.Next()
	request.Scan(&id, &name, &value, &description)

	var dataOut = map[string]string{
		"id":          strconv.Itoa(id),
		"name":        name,
		"value":       strconv.Itoa(value),
		"description": description,
	}

	return c.JSON(http.StatusOK, dataOut)
}

func CheckFlag(c echo.Context) error {
	var flag = c.FormValue("flag")
	var taskId = c.FormValue("taskid")

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*common.JwtCustomClaims)
	userid := claims.UserId

	var userCommandId int
	var getCommandId, _ = database.DB.Query("SELECT commandid FROM users WHERE id=?", userid)
	getCommandId.Next()
	getCommandId.Scan(&userCommandId)

	var getRightAnswer, errGetTaskAnswer = database.DB.Query("SELECT flag, value, contestid, type FROM tasks left join contests on tasks.contestid = contests.id where tasks.id=?", taskId)

	if errGetTaskAnswer != nil {
		return c.JSON(http.StatusServiceUnavailable, errGetTaskAnswer)
	}

	// Right answer for check and point if a true flag was submitted
	var (
		rightAnswer, contestType string
		points, contestid        int
	)
	getRightAnswer.Next()
	getRightAnswer.Scan(&rightAnswer, &points, &contestid, &contestType)

	var pwnedStatus int
	var (
		checkPwned    *sql.Rows
		errCheckPwned error
	)

	//log.Println(contestType)

	if contestType == "alone" {
		checkPwned, errCheckPwned = database.DB.Query("SELECT COUNT(*) FROM pwnedby WHERE userid=? AND taskid=?", userid, taskId)
		if errCheckPwned != nil {
			return c.JSON(http.StatusServiceUnavailable, errCheckPwned)
		}
		checkPwned.Next()
		checkPwned.Scan(&pwnedStatus)
		if pwnedStatus == 1 {
			return c.JSON(http.StatusOK, map[string]string{
				"result": "already",
			})
		}
		database.DB.Query("INSERT INTO pwnedby (userid, taskid) VALUES(?,?)", userid, taskId) // Set task as accepted by this user

	} else if contestType == "team" {
		//log.Println(userCommandId, taskId)
		checkPwned, errCheckPwned = database.DB.Query("SELECT COUNT(*) FROM pwnedby WHERE command_id=? AND taskid=?", userCommandId, taskId)
		if errCheckPwned != nil {
			return c.JSON(http.StatusServiceUnavailable, errCheckPwned)
		}
		checkPwned.Next()
		checkPwned.Scan(&pwnedStatus)
		if pwnedStatus == 1 {
			return c.JSON(http.StatusOK, map[string]string{
				"result": "already",
			})
		}

		database.DB.Query("INSERT INTO pwnedby (userid, taskid, command_id) VALUES(?,?, ?)", userid, taskId, userCommandId) // Set task as accepted by this user
	}

	if rightAnswer == flag {
		//database.DB.Query("UPDATE users SET `flagright`=`flagright`+1, `points`=`points`+? WHERE id=?", points, userid) // Increment for bruteforcer factor and +points

		var checkExistInRatingTable, errCheckExist = database.DB.Query("SELECT COUNT(*) FROM rating WHERE contest_id = ? and userid = ?", contestid, userid);
		if errCheckExist != nil {
			log.Println(errCheckExist)
		}
		var existStatus int
		checkExistInRatingTable.Next()
		checkExistInRatingTable.Scan(&existStatus)

		return c.JSON(http.StatusOK, map[string]bool{
			"result": true,
		})
	} else {
		database.DB.Query("UPDATE users SET `flagfalse`=`flagfalse`+1 WHERE id=?", userid) // Increment for bruteforcer factor

		return c.JSON(http.StatusOK, map[string]bool{
			"result": false,
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"id":     strconv.Itoa(userid),
		"flag":   flag,
		"taskid": taskId,
		"answer": rightAnswer,
	})
}
