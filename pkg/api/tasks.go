package api

import (
	"net/http"

	"goFinalProject/pkg/db"
)

const tasksLimit = 50

type TasksResponse struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w, http.MethodGet)
		return
	}

	tasks, err := db.Tasks(tasksLimit, r.URL.Query().Get("search"))
	if err != nil {
		writeInternalError(w, "не удалось получить список задач", err)
		return
	}

	writeJSON(w, TasksResponse{Tasks: tasks}, http.StatusOK)
}
