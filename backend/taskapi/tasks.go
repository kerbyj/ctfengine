package taskapi

import (
	"github.com/kerbyj/ctfengine/backend/common"
	"github.com/kerbyj/ctfengine/backend/database"
	"database/sql"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"runtime"
	"time"
)

type contestStruct struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	TasksCount int    `json:"tasks_count"`
}

func GetContestList(c echo.Context) error {
	var getAllContests, errGetAllContests = database.DB.Query("select id, name, type, (select count(*) from tasks where contestid=contests.id) as taskscount from contests where visibility=true")
	defer getAllContests.Close()

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
	defer getCommandId.Close()

	getCommandId.Next()
	getCommandId.Scan(&userCommandId)

	var GetAllTasks, errGetAliveTasks = database.DB.Query("SELECT tasks.id, category, tasks.name, value, contests.name FROM tasks left join contests on tasks.contestid = contests.id where contests.visibility = true order by contestid")
	defer GetAllTasks.Close()

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
	defer getAlreadySolvedTasks.Close()

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
	defer request.Close()

	if errGetTask != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	var value, id int
	var name, description, category, contest string

	type Attachment struct {
		Id int `json:"id"`
		Name string `json:"name"`
	}

	type TaskInFoOut struct {
		Id int `json:"id"`
		Name string `json:"name"`
		Value int `json:"value"`
		Description string `json:"description"`
		Category string `json:"category"`
		Contest string `json:"contest"`
		Attachments []Attachment `json:"attachments"`
	}

	request.Next()
	request.Scan(&id, &name, &value, &description, &category, &contest)
	var getAttachmentsForTask, errGetAttachmentsForTask = database.DB.Query("select id, name  from attachments where taskid=?", requestedTask)
	defer getAttachmentsForTask.Close()

	var(
		allAttachments []Attachment
		attachmentId int
		attachmentName string
	)

	if errGetAttachmentsForTask != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	for getAttachmentsForTask.Next(){
		getAttachmentsForTask.Scan(&attachmentId, &attachmentName)
		allAttachments = append(allAttachments, Attachment{
			attachmentId,
			attachmentName,
		})
	}

	var dataOut = TaskInFoOut{
		id,
		name,
		value,
		description,
		category,
		contest,
		allAttachments,
	}

	return c.JSON(http.StatusOK, dataOut)
}

/*
func GetFileById(c echo.Context) error {
	var requestedTask = c.Param("id")
	var request, errGetTask = database.DB.Query("select name from attachments", requestedTask)

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
*/

func CheckFlag(c echo.Context) error {
	var flag = c.FormValue("flag")
	var taskId = c.FormValue("taskid")

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*common.JwtCustomClaims)
	userid := claims.UserId

	var userCommandId int
	var username string

	var getCommandId, _ = database.DB.Query("SELECT username, commandid FROM users WHERE id=?", userid)
	defer getCommandId.Close()

	getCommandId.Next()
	getCommandId.Scan(&username, &userCommandId)



	var getRightAnswer, errGetTaskAnswer = database.DB.Query("SELECT flag, value, contestid, type, tasks.name, contests.permit FROM tasks left join contests on tasks.contestid = contests.id where tasks.id=?", taskId)
	defer getRightAnswer.Close()

	if errGetTaskAnswer != nil {
		log.Println(errGetTaskAnswer)
		return c.JSON(http.StatusServiceUnavailable, errGetTaskAnswer)
	}

	// Right answer for check and point if a true flag was submitted
	var (
		rightAnswer, contestType, taskName string
		points, contestid, permit        int
	)
	getRightAnswer.Next()
	getRightAnswer.Scan(&rightAnswer, &points, &contestid, &contestType, &taskName, &permit)

	if permit == 0 {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "error",
			"error": "time over",
		})
	}

	var pwnedStatus int
	var (
		checkPwned    *sql.Rows
		errCheckPwned error
	)

	//log.Println(contestType)

	if userCommandId == 0 && contestType == "team" {
		log.Println("you need to create command")
		return c.JSON(http.StatusOK, map[string]string{
			"result": "you need to create command",
		})
	}

	if contestType == "alone" {
		checkPwned, errCheckPwned = database.DB.Query("SELECT COUNT(*) FROM pwnedby WHERE userid=? AND taskid=?", userid, taskId)
		defer checkPwned.Close()

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


	} else if contestType == "team" {
		checkPwned, errCheckPwned = database.DB.Query("SELECT COUNT(*) FROM pwnedby WHERE command_id=? AND taskid=?", userCommandId, taskId)
		defer checkPwned.Close()

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
	}

	if rightAnswer == flag {
		if contestType == "alone" {
			var checkExistInRatingTable, errCheckExist = database.DB.Query("SELECT COUNT(*) FROM rating WHERE contest_id = ? and userid = ? ", contestid, userid)
			defer checkExistInRatingTable.Close()

			if errCheckExist != nil {
				log.Println(errCheckExist)
			}

			var existStatus int
			checkExistInRatingTable.Next()
			checkExistInRatingTable.Scan(&existStatus)

			if existStatus == 0 {
				var cc, _ = database.DB.Query("INSERT INTO rating (contest_id, userid, points) values (?, ?, ?)", contestid, userid, points)
				defer cc.Close()
			} else {
				var cc, _ = database.DB.Query("UPDATE rating SET points = points + ? WHERE contest_id=? AND userid=?", points, contestid, userid)
				defer cc.Close()
			}
			var cc, _ = database.DB.Query("UPDATE users SET `flagright`=`flagright`+1 WHERE id=?", userid)
			defer cc.Close()

			var cc2, _ = database.DB.Query("INSERT INTO pwnedby (userid, taskid, contestid, time) VALUES(?,?,?,?)", userid, taskId, contestid,time.Now()) // Set task as accepted by this user
			defer cc2.Close()

		} else if contestType == "team" {
			var commandName string
			var getCommandName, _ = database.DB.Query("SELECT name from command where id=?", userCommandId)
			getCommandName.Next()
			getCommandName.Scan(&commandName)

			log.Println(commandName, "sent the correct flag(", flag, ") for task", taskId, taskName, "in", time.Now())

			var checkExistInRatingTable, errCheckExist = database.DB.Query("SELECT COUNT(*) FROM rating WHERE contest_id = ? and team = ? ", contestid, userCommandId)
			defer checkExistInRatingTable.Close()

			if errCheckExist != nil {
				log.Println(errCheckExist)
			}
			var existStatus int
			checkExistInRatingTable.Next()
			checkExistInRatingTable.Scan(&existStatus)

			if existStatus == 0 {
				var cc, _ = database.DB.Query("INSERT INTO rating (contest_id, team, points) values (?, ?, ?)", contestid, userCommandId, points)
				defer cc.Close()

			} else {
				var cc, _ = database.DB.Query("UPDATE rating SET points = points + ? WHERE contest_id=? AND team=?", points, contestid, userCommandId)
				defer cc.Close()
			}

			var cc, _ = database.DB.Query("INSERT INTO pwnedby (userid, taskid, command_id, contestid, time) VALUES(?,?,?,?,?)", userid, taskId, userCommandId, contestid, time.Now()) // Set task as accepted by this user
			defer cc.Close()
		}

		return c.JSON(http.StatusOK, map[string]bool{
			"result": true,
		})
	} else {
		var cc, _ = database.DB.Query("UPDATE users SET `flagfalse`=`flagfalse`+1 WHERE id=?", userid) // Increment for bruteforcer factor
		cc.Close()

		return c.JSON(http.StatusOK, map[string]bool{
			"result": false,
		})
	}
}
