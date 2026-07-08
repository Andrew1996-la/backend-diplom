package api

import (
	"net/http"

	"github.com/Andrew1996-la/backend-diplom/pkg/db"
)

type tasksResponse struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(50)
	if err != nil {
		writeJson(w, errorResponse{Error: err.Error()})
		return
	}

	writeJson(w, tasksResponse{
		Tasks: tasks,
	})
}
