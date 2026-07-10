package api

import (
	"net/http"

	"github.com/Andrew1996-la/backend-diplom/pkg/db"
)

type tasksResponse struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("search")

	tasks, err := db.Tasks(50, search)
	if err != nil {
		writeJson(w, errorResponse{
			Error: err.Error(),
		})
		return
	}

	writeJson(w, tasksResponse{
		Tasks: tasks,
	})
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeJson(w, map[string]string{
			"error": "Не указан идентификатор",
		})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, errorResponse{
			Error: err.Error(),
		})
		return
	}

	writeJson(w, task)
}
