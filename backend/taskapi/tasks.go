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
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	TasksCount int    `json:"tasks_count"`
}

func GetContestList(c echo.Context) error {
	var getAllContests, errGetAllContests = database.DB.Query("select id, name, type, (select count(*) from tasks where contestid=contests.id) as taskscount from contests where visibility=true")
	if errGetAllContests != nil {
		return c.String(http.StatusBadRequest, "shit happens")
	}

	var (
		dataOut           []contestStruct
		id, taskscount    int
		name, contestType string
	)

	for getAllContests.Next() {
		getAllContests.Scan(&id, &name, &contestType, &taskscount)
		dataOut = append(dataOut, contestStruct{
			Id:         id,
			Name:       name,
			Type:       contestType,
			TasksCount: taskscount,
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

	var GetAllTasks, errGetAliveTasks = database.DB.Query("SELECT tasks.id, category, tasks.name, value, contests.name FROM tasks left join contests on tasks.contestid = contests.id where contests.visibility = true order by contestid")
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
	var request, errGetTask = database.DB.Query("SELECT tasks.id, tasks.name, value, description, category, contests.name FROM tasks left join contests on tasks.contestid = contests.id WHERE tasks.id=?", requestedTask)

	if errGetTask != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	var value, id int
	var name, description, category, contest string

	request.Next()
	request.Scan(&id, &name, &value, &description, &category, &contest)

	var dataOut = map[string]string{
		"id":          strconv.Itoa(id),
		"name":        name,
		"value":       strconv.Itoa(value),
		"description": description,
		"category":    category,
		"contest":     contest,
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
		database.DB.Query("INSERT INTO pwnedby (userid, taskid, contestid) VALUES(?,?,?)", userid, taskId, contestid) // Set task as accepted by this user

	} else if contestType == "team" {
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

		database.DB.Query("INSERT INTO pwnedby (userid, taskid, command_id, contestid) VALUES(?,?,?,?)", userid, taskId, userCommandId, contestid) // Set task as accepted by this user
	}

	if rightAnswer == flag {

		if contestType == "alone" {

			var checkExistInRatingTable, errCheckExist = database.DB.Query("SELECT COUNT(*) FROM rating WHERE contest_id = ? and userid = ? ", contestid, userid)
			if errCheckExist != nil {
				log.Println(errCheckExist)
			}
			var existStatus int
			checkExistInRatingTable.Next()
			checkExistInRatingTable.Scan(&existStatus)

			if existStatus == 0 {
				database.DB.Query("INSERT INTO rating (contest_id, userid, points) values (?, ?, ?)", contestid, userid, points)
			} else {
				database.DB.Query("UPDATE rating SET points = points + ? WHERE contest_id=? AND userid=?", points, contestid, userid)
			}

		} else if contestType == "team" {
			var checkExistInRatingTable, errCheckExist = database.DB.Query("SELECT COUNT(*) FROM rating WHERE contest_id = ? and team = ? ", contestid, userCommandId)
			if errCheckExist != nil {
				log.Println(errCheckExist)
			}
			var existStatus int
			checkExistInRatingTable.Next()
			checkExistInRatingTable.Scan(&existStatus)

			if existStatus == 0 {
				database.DB.Query("INSERT INTO rating (contest_id, team, points) values (?, ?, ?)", contestid, userCommandId, points)
			} else {
				database.DB.Query("UPDATE rating SET points = points + ? WHERE contest_id=? AND team=?", points, contestid, userCommandId)
			}
		}

		return c.JSON(http.StatusOK, map[string]bool{
			"result": true,
		})
	} else {
		database.DB.Query("UPDATE users SET `flagfalse`=`flagfalse`+1 WHERE id=?", userid) // Increment for bruteforcer factor

		return c.JSON(http.StatusOK, map[string]bool{
			"result": false,
		})
	}
}
