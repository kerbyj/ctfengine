package boardapi

import (
	"ctfEngine/backend/database"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"strconv"
)

func BoardStats(c echo.Context) error {

	var request, errorRequestBoardStat = database.DB.Query("select " +
		"(select count(*) from ctfengine.tasks)," +
		"(select count(*) from ctfengine.users)," +
		"(select username from ctfengine.users order by points desc limit 1)")

	if errorRequestBoardStat != nil {
		log.Println(errorRequestBoardStat)
	}
	defer request.Close()

	var tasksCount, userCount int
	var topUserName string

	request.Next()
	request.Scan(&tasksCount, &userCount, &topUserName)


	var taskDataRequest, errorTaskDataRequest = database.DB.Query("SELECT category, count(category) " +
		"from tasks " +
		"group by category")

	if errorTaskDataRequest != nil {
		log.Println("ERROR!!!!", errorRequestBoardStat)
	}
	defer taskDataRequest.Close()


	var dataOut = map[string]map[string]string{
		"tasksStats": {
			"tasks_count":          strconv.Itoa(tasksCount),
		},
		"boardStats":{
			"user_count":       strconv.Itoa(userCount),
			"TOP 1 pwner":        topUserName,
		},
	}

	var category string
	var categoryCount int
	for taskDataRequest.Next() {
		taskDataRequest.Scan(&category, &categoryCount)
		dataOut["tasksStats"][category] = strconv.Itoa(categoryCount)
	}

	return c.JSON(http.StatusOK, dataOut)
}
