package taskapi

import (
	"ctfEngine/backend/database"
	"github.com/labstack/echo"
	"net/http"
)

type taskStruct struct {
	Category    string `json:"Category"`
	Name        string `json:"Name"`
	Description string `json:"description"`
	Value       int    `json:"value"`
	Flag        string `json:"flag"`
}

func GetAlwaysAliveTasks(c echo.Context) error {
	var request, errGetAliveTasks = database.DB.Query("SELECT * FROM tasks ORDER BY category")

	if errGetAliveTasks != nil {
		return c.String(http.StatusBadRequest, "shit happens")
	}

	var (
		taskOut                           = make(map[int]taskStruct)
		id, value                         int
		category, name, description, flag string
	)
	for request.Next() {
		request.Scan(&id, &category, &name, &description, &value, &flag)

		taskOut[id] = taskStruct{
			Category:    category,
			Name:        name,
			Description: description,
			Value:       value,
			Flag:        flag,
		}
	}
	return c.JSON(http.StatusOK, taskOut)
}
