package taskapi

import (
	"ctfEngine/backend/database"
	"github.com/labstack/echo"
	"net/http"
	"runtime"
)

type taskStruct struct {
	Name        string `json:"Name"`
	Description string `json:"description"`
	Value       int    `json:"value"`
	Flag        string `json:"flag"`
	Status      string `json:"status"` // Solved or not
}

func GetAlwaysAliveTasks(c echo.Context) error {
	var request, errGetAliveTasks = database.DB.Query("SELECT * FROM tasks ORDER BY category")

	if errGetAliveTasks != nil {
		return c.String(http.StatusBadRequest, "shit happens")
	}

	var (
		taskOut                           = make(map[string][]taskStruct)
		id, value                         int
		category, name, description, flag string
	)

	for request.Next() {
		request.Scan(&id, &category, &name, &description, &value, &flag)

		taskOut[category] = append(taskOut[category], taskStruct{
			Name:        name,
			Description: description,
			Value:       value,
			Flag:        flag,
		})
		// Please, create issue if you have better way for this method
	}

	runtime.GC()
	return c.JSON(http.StatusOK, taskOut)
}
