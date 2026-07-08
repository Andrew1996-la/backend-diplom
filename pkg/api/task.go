package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Andrew1996-la/backend-diplom/pkg/db"
)

type idResponse struct {
	ID string `json:"id"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func writeJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func checkDate(task *db.Task) error {
	now := time.Now()
	today := now.Format(dateFormat)

	if task.Date == "" {
		task.Date = today
	}

	_, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return err
	}

	if task.Date < today {
		if task.Repeat == "" {
			task.Date = today
		} else {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return err
			}
			task.Date = next
		}
	}

	return nil
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJson(w, errorResponse{Error: err.Error()})
		return
	}

	if task.Title == "" {
		writeJson(w, errorResponse{Error: "task title is required"})
		return
	}

	if err := checkDate(&task); err != nil {
		writeJson(w, errorResponse{Error: err.Error()})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJson(w, errorResponse{Error: err.Error()})
		return
	}

	writeJson(w, idResponse{
		ID: strconv.FormatInt(id, 10),
	})
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	default:
		writeJson(w, errorResponse{Error: "method not allowed"})
	}
}
