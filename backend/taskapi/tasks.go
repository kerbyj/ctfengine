package taskapi

import (
	"ctfEngine/backend/database"
	"github.com/labstack/echo"
	"net/http"
	"runtime"
	"strconv"
)

type taskStruct struct {
	Name  string `json:"Name"`
	Value int    `json:"value"`
	Id    int    `json:"id"`
}

func GetAlwaysAliveTasks(c echo.Context) error {
	var request, errGetAliveTasks = database.DB.Query("SELECT id, category, name, value FROM tasks ORDER BY category")

	if errGetAliveTasks != nil {
		return c.String(http.StatusBadRequest, "shit happens")
	}

	var (
		taskOut        = make(map[string][]taskStruct)
		id, value      int
		category, name string
	)

	for request.Next() {
		request.Scan(&id, &category, &name, &value)

		taskOut[category] = append(taskOut[category], taskStruct{
			Name: name,
			Value: value,
			Id: id,
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
		"id": strconv.Itoa(id),
		"name": name,
		"value": strconv.Itoa(value),
		"description": description,
	}

	return c.JSON(http.StatusOK, dataOut)
}
